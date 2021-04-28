package main

import (
	"context"
	"log"

	"github.com/pkg/errors"

	"github.com/WiggiLi/file-sharing-api/service"
	"github.com/WiggiLi/file-sharing-api/store"
	"github.com/WiggiLi/file-sharing-api/controller"
	"github.com/WiggiLi/file-sharing-api/lib/logger"
)

func run(errc chan<- error) {
	ctx := context.Background()

	l := logger.Get()

	// Init repository store
	store, err := store.New(ctx)
	if err != nil {
		errc <- errors.Wrap(err, "store.New failed")
		return
	}

	// Init service manager
	serviceManager, err := service.NewManager(ctx, store)
	if err != nil {
		errc <- errors.Wrap(err, "manager.New failed")  		
		return
	}

	controller := controller.NewController(ctx, serviceManager, l)
	controller.Start(errc)
}



func main() {
	log.Print("Server is preparing to start...")

	errc := make(chan error)
	go run(errc)
	if err := <-errc; err != nil {
		log.Fatal(err)
	}
}