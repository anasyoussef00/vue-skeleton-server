package main

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"github.com/youssef-182/vue-skeleton-server/pkg/db"
	"github.com/youssef-182/vue-skeleton-server/pkg/router"
	"log"
	"net/http"
	"os"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("An error has occurred while trying to load the .env file: %v", err)
	}

	if err := db.Connect("mysql", os.Getenv("DB_DSN")); err != nil {
		log.Fatalf("An error has occurred while trying to connect to database: %v", err)
	}

	if err := db.Migrate("C:\\Users\\Anas Youssef\\Documents\\projects\\vue-skeleton-server\\sql"); err != nil {
		log.Fatal(err)
	}
}

func main() {
	r := router.Init()

	if err := http.ListenAndServe(":1337", r); err != nil {
		log.Fatalf("There was an error listening to port :1337. %v", err)
	}
}
