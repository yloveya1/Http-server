package main

import (
	"context"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"go.mod/internal/config"
	"go.mod/pkg/handler"
	"go.mod/pkg/repository"
	"go.mod/pkg/service"
	"log"
	"net"
	"net/http"
	"time"
)

func main() {
	router := httprouter.New()
	cfg := config.GetConf()
	postgreSQLClient, _ := repository.NewClient(context.TODO(), cfg)
	serv := service.NewService(postgreSQLClient)
	handler := handler.NewHandler(serv)
	handler.Register(router)

	Start(router, cfg)

}

func Start(router *httprouter.Router, conf *config.Config) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", conf.Port))
	if err != nil {
		panic(err)
	}
	server := &http.Server{
		Handler:      router,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatalln(server.Serve(listener))
}
