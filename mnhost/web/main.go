package main

import (
	"net/http"

	"github.com/micro/go-micro/util/log"

	"github.com/micro/go-micro/web"
	"vircle.io/mnhost/web/handler"
)

func main() {
	// create new web service
	service := web.NewService(
		web.Name("go.micro.web.web"),
		web.Version("latest"),
		web.Address("0.0.0.0:8888"),
	)

	// initialise service
	if err := service.Init(); err != nil {
		log.Fatal(err)
	}

	// register html handler
	service.Handle("/", http.FileServer(http.Dir("html")))

	// register call handler
	service.HandleFunc("/web/call", handler.WebCall)

	// run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
