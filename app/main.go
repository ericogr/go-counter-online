package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-kit/kit/log"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/viper"
)

var (
	PATH_ROOT        = "/"
	PATH_COUNT_POST  = "/count/{uuid:[0-9A-F]{8}-[0-9A-F]{4}-[5][0-9A-F]{3}-[89AB][0-9A-F]{3}-[0-9A-F]{12}}/{name:[a-zA-Z][a-zA-Z0-9]{0,15}}"
	PATH_COUNT_GET   = "/count/{uuid:[0-9A-F]{8}-[0-9A-F]{4}-[5][0-9A-F]{3}-[89AB][0-9A-F]{3}-[0-9A-F]{12}}"
	PATH_METRICS_GET = "/metrics"
)

func init() {
	viper.AutomaticEnv()
	viper.SetEnvPrefix("COUNTER")

	viper.SetDefault("Port", 8080)
	viper.SetDefault("DatabaseConfiguration", "")
	viper.SetDefault("HideDatabaseConfigurationOutput", false)

	if viper.IsSet("Config") {
		viper.SetConfigFile(viper.GetString("Config"))
	} else {
		viper.SetConfigName("config")
		viper.AddConfigPath(".")
	}
	viper.SetConfigType("json")
	viper.SetConfigType("yaml")

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	fmt.Println("Configuration:")
	fmt.Printf("- port: %d\n", viper.GetInt("Port"))
	fmt.Printf("- hide-database-configuration-output: %t\n", viper.GetBool("HideDatabaseConfigurationOutput"))
	if !viper.GetBool("HideDatabaseConfigurationOutput") {
		fmt.Printf("- database-configuration: %v\n", viper.GetString("DatabaseConfiguration"))
	}
}

func main() {
	var logger log.Logger = log.NewLogfmtLogger(os.Stderr)
	logger = log.With(logger, "listen", viper.GetInt("Port"), "caller", log.DefaultCaller)

	fieldsKeys := []string{"method", "error"}
	requestCount := kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
		Subsystem: "counter_service",
		Name:      "request_count",
		Help:      "Number of requests received.",
	}, fieldsKeys)
	requestLatency := kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
		Subsystem: "counter_service",
		Name:      "request_latency_microseconds",
		Help:      "Total duration of requests in microseconds.",
	}, fieldsKeys)

	var databaseParams string = viper.GetString("DatabaseConfiguration")
	var svc CounterService = &PostgresCounterService{
		DatabaseParams: databaseParams,
	}
	svc = DefaultLoggingMiddleware(logger)(svc)
	svc = DefaultInstrumentingMiddleware(requestCount, requestLatency)(svc)
	err := svc.Init()
	if err != nil {
		panic(err)
	}

	counterCreateHandler := httptransport.NewServer(
		makeCountCreateEndpoint(svc),
		decodeCountCreateRequest,
		encodeResponse,
	)

	counterIncrementHandler := httptransport.NewServer(
		makeCountIncrementEndpoint(svc),
		decodeCountIncrementRequest,
		encodeResponse,
	)

	rootHandler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Hi there!\n")
		fmt.Fprintf(w, "- call GET /metrics to get metrics in Prometheus format\n")
		fmt.Fprintf(w, "- call GET /count/{uuidv5} to validate uuid\n")
		fmt.Fprintf(w, "- call POST /count/{uuidv5}/{name} to create new counter\n")
	}

	router := mux.NewRouter()
	router.Handle(PATH_COUNT_POST, counterCreateHandler).Methods("POST")
	router.Handle(PATH_COUNT_GET, counterIncrementHandler).Methods("GET")
	router.Handle(PATH_METRICS_GET, promhttp.Handler()).Methods("GET")
	router.HandleFunc(PATH_ROOT, rootHandler).Methods("GET")

	addr := fmt.Sprintf(":%d", viper.GetInt("Port"))
	srv := &http.Server{
		Handler:      router,
		Addr:         addr,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	go func() {
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			logger.Log("err", err)
		}
	}()

	waitForShutdown(svc)
}

func waitForShutdown(svc CounterService) {
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	signalx := <-s

	fmt.Println("Shutting down:", signalx)

	err := svc.Terminate()
	if err != nil {
		fmt.Printf("error while terminating service: %s", err)
	} else {
		fmt.Println("service terminated with success")
	}
}
