package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ericogr/go-counter-online/rest"
	"github.com/ericogr/go-counter-online/services"
	"github.com/ericogr/go-counter-online/storage"
	"github.com/gorilla/mux"
	"github.com/spf13/viper"
)

var (
	PATH_ROOT       = "/"
	PATH_COUNT_GET  = "/count/{uuid:[0-9A-F]{8}-[0-9A-F]{4}-[5][0-9A-F]{3}-[89AB][0-9A-F]{3}-[0-9A-F]{12}}"
	PATH_COUNT_POST = "/count/{uuid:[0-9A-F]{8}-[0-9A-F]{4}-[5][0-9A-F]{3}-[89AB][0-9A-F]{3}-[0-9A-F]{12}}/{name:[a-zA-Z][a-zA-Z0-9]{0,15}}"
)

func init() {
	fmt.Printf("Initializing...\n")

	viper.AutomaticEnv()
	viper.SetEnvPrefix("COUNTER")

	viper.SetDefault("Port", 8080)
	viper.SetDefault("Database", "memory")
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
	fmt.Printf("- database: %v\n", viper.GetString("Database"))
	fmt.Printf("- hide-database-configuration-output: %t\n", viper.GetBool("HideDatabaseConfigurationOutput"))
	if !viper.GetBool("HideDatabaseConfigurationOutput") {
		fmt.Printf("- database-configuration: %v\n", viper.GetString("DatabaseConfiguration"))
	}
}

func waitForShutdown(dataStorage storage.CounterData) {
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	signalx := <-s

	fmt.Println("Shutting down:", signalx)

	err := dataStorage.Terminate()
	if err != nil {
		fmt.Printf("error while terminating database: %s", err)
	}
}

func getStorageInstance() (storage.CounterData, error) {
	return storage.GetStoreInstance(
		viper.GetString("Database"),
		viper.GetString("DatabaseConfiguration"),
	)
}

func doRoute() {
	router := mux.NewRouter()
	storageInstance, err := getStorageInstance()
	if err != nil {
		panic(err)
	}

	countRest := rest.CounterRest{
		CounterService: &services.DefaultCounterService{
			CounterData: storageInstance,
		},
	}

	router.HandleFunc(PATH_ROOT, countRest.DoNothing)
	router.HandleFunc(PATH_COUNT_GET, countRest.GetCount).Methods("GET")
	router.HandleFunc(PATH_COUNT_POST, countRest.CreateCount).Methods("POST")

	addr := fmt.Sprintf(":%d", viper.GetInt("Port"))
	srv := &http.Server{
		Handler:      router,
		Addr:         addr,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	go func() {
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe: %v", err)
		}
	}()

	waitForShutdown(storageInstance)
}

func main() {
	fmt.Printf("Listening on port %d\n", viper.GetInt("Port"))
	doRoute()
}
