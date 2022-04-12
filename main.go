package main

import (
	"flag"
	"github.com/joho/godotenv"
	"github.com/romankarpowich/ozon/app/router"
	"github.com/romankarpowich/ozon/repository"
	"log"
	"net/http"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalln(err)
	}

	flag.StringVar(repository.MemoryType, "memory", "postgresql", "Укажите один из двух типов хранилища - postgresql или in-memory.")
	flag.Parse()

	repository.InitStore()
	router.InitRouter()

	log.Println("App starts on port :8080 and memory - " + *repository.MemoryType)

	err = http.ListenAndServe(":8080", router.Mux)
	if err != nil {
		log.Fatalln(err)
	}
}
