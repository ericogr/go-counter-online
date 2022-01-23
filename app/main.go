package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/ericogr/go-counter-online/routes"
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

func main() {
	fmt.Printf("Listening on port %d\n", viper.GetInt("Port"))
	doRoute()
}

func doRoute() {
	r := mux.NewRouter()
	r.HandleFunc(PATH_ROOT, routes.DoNothingRoute)
	r.HandleFunc(PATH_COUNT_GET, routes.GetCountRoute).Methods("GET")
	r.HandleFunc(PATH_COUNT_POST, routes.CreateCountRoute).Methods("POST")
	log.Fatal(
		http.ListenAndServe(
			fmt.Sprintf(":%d", viper.GetInt("Port")),
			r,
		),
	)
}
