package main

import (
	"gobank/routes"
	"gobank/storage"
	"log"
)

func main() {
	//create Storeage

	store, err := storage.NewPostgresStore()
	if err != nil {
		log.Fatal(err)
	}

	//init storage
	if err := store.Init(); err != nil {
		log.Fatal(err)
	}

	//init Apiserver
	server := routes.NewServer(":2000", store)

	//run the apiserver
	server.Run()

}
