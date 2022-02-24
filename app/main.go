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
	"github.com/go-kit/log/level"
	"github.com/gorilla/mux"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var (
	logger                      log.Logger
	PATH_ROOT                   = "/"
	PATH_COUNT_POST             = "/count/{uuid:[0-9A-F]{8}-[0-9A-F]{4}-[5][0-9A-F]{3}-[89AB][0-9A-F]{3}-[0-9A-F]{12}}/{name:[a-zA-Z][a-zA-Z0-9]{0,15}}"
	PATH_COUNT_GET              = "/count/{uuid:[0-9A-F]{8}-[0-9A-F]{4}-[5][0-9A-F]{3}-[89AB][0-9A-F]{3}-[0-9A-F]{12}}"
	PATH_METRICS_GET            = "/metrics"
	FLAG_PORT                   = "port"
	FLAG_DEBUG                  = "debug"
	FLAG_DATABASE_CONFIGURATION = "databaseConfiguration"
)

func init() {
	pflag.Int(FLAG_PORT, 8080, "Port to listen")
	pflag.String(FLAG_DATABASE_CONFIGURATION, "", "Database configuration or connection string")
	pflag.String(FLAG_DEBUG, "info", "Enable debug mode")
	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)

	viper.AutomaticEnv()
	viper.SetEnvPrefix("COUNTER")

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

	logger = log.NewLogfmtLogger(os.Stderr)
	logger = level.NewFilter(
		logger,
		parseLogOption(viper.GetString(FLAG_DEBUG)),
	)
	logger = log.With(
		logger,
		"ts",
		log.DefaultTimestampUTC,
		"caller",
		log.DefaultCaller,
	)

	printSettings()
}

func main() {
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

	var databaseParams string = viper.GetString(FLAG_DATABASE_CONFIGURATION)
	var svc CounterService = &PostgresCounterService{
		DatabaseParams: databaseParams,
	}
	svc = DefaultLoggingMiddleware(logger, viper.GetBool(FLAG_DEBUG))(svc)
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

	addr := fmt.Sprintf(":%d", viper.GetInt(FLAG_PORT))
	srv := &http.Server{
		Handler:      router,
		Addr:         addr,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	go func() {
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			level.Error(logger).Log("error", err)
		}
	}()

	waitForShutdown(svc)
}

func waitForShutdown(svc CounterService) {
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	signalx := <-s

	logger.Log("Shutting down:", signalx)

	err := svc.Terminate()
	if err != nil {
		level.Error(logger).Log("error", fmt.Sprintf("error while terminating service: %s", err))
	}
}

func printSettings() {
	for k, v := range viper.AllSettings() {
		level.Debug(logger).Log("setting", k, "value", v)
	}
}

func parseLogOption(logLevel string) level.Option {
	switch logLevel {
	case "debug":
		return level.AllowDebug()
	case "info":
		return level.AllowInfo()
	case "warn":
		return level.AllowWarn()
	case "error":
		return level.AllowError()
	case "fatal":
		return level.AllowError()
	default:
		return level.AllowInfo()
	}
}
