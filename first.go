package dblayer

import (
	"awesomeProject/configuration"
	"awesomeProject/eventservice/rest"
	"awesomeProject/peresistence/dblayer"
	"flag"
	"fmt"
	"log"

	//_ "github.com/go-sql-driver/mysql"
	//"github.com/gorilla/mux"
	//"net/http"

	//"github.com/jmoiron/sqlx"
)

func main() {
	confPath := flag.String("conf", `.\configuration\config.json`, "flag to set the path to the configuration json file")
	flag.Parse()
	//extract configuration
	config, _ := configuration.ExtractConfiguration(*confPath)
	fmt.Println("Connecting to database")
	dbhandler, _ := dblayer.NewPersistenceLayer(config.Databasetype, config.DBConnection)
	//RESTful API start
	httpErrChan, httptlsErrChan := rest.ServeAPI(config.RestfulEndpoint, config.RestfulTLSEndPint, dbhandler)
	select {
	case err := <-httpErrChan:
		log.Fatal("HTTP Error: ", err)
	case err := <-httptlsErrChan:
		log.Fatal("HTTPS Error: ", err)
	}
}
