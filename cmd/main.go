package main

import (
	"context"
	"document-service/config"
	"document-service/handler"
	"document-service/models"
	"document-service/services"
	"document-service/store"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	db, err := store.NewStore(ctx, config.Config.DSN())
	if err != nil {
		log.Fatalf("не удалось подключиться к базе данных: %v", err)
	}
	defer db.Close()

	services.DB = db
	services.Session = make(map[string]struct{})
	services.Documents = make(map[uint64]models.Document)

	router := mux.NewRouter()
	handlers.InitRouter(router)

	log.Println("Запуск веб-сервера на http://127.0.0.1:8080")
	log.Fatal(http.ListenAndServe(":"+config.Config.Port, router))
}
