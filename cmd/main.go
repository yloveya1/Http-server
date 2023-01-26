package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/nats-io/stan.go"
	"go.mod/internal/cache"
	"go.mod/internal/config"
	"go.mod/internal/domain"
	"go.mod/pkg/handler"
	"go.mod/pkg/repository"
	"go.mod/pkg/service"
	"io"
	"log"
	"math/rand"
	"net"
	"net/http"
	"time"
)

func main() {
	log.Println("Router is created")
	router := httprouter.New()
	log.Println("Config is initialized")
	cfg := config.GetConf()
	postgreSQLClient, err := repository.NewClient(context.Background(), cfg)
	if err != nil {
		log.Fatalln(err)
	}
	serv := service.NewService(postgreSQLClient, cache.New())
	err = serv.FindAllOrders(context.Background())
	if err != nil {
		log.Fatalln(err)
	}
	handler := handler.NewHandler(serv)
	handler.Register(router)
	log.Println("Handler is registered")
	go startNats(serv, cfg)
	Start(router, cfg)
}

func Start(router *httprouter.Router, conf *config.Config) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", conf.Port))
	log.Println("Listen TCP")
	if err != nil {
		log.Fatalln(err)
	}
	server := &http.Server{
		Handler:      router,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatalln(server.Serve(listener))
}

func logCloser(c io.Closer) {
	if err := c.Close(); err != nil {
		log.Fatalln(err)
	}
}

func SubscriberNats(s repository.Repository, conn stan.Conn) {

	var err error

	_, err = conn.Subscribe("NewOrder", func(msg *stan.Msg) {

		var ord domain.Order
		if err = json.Unmarshal(msg.Data, &ord); err != nil {
			log.Fatalln(err)
		}
		s.AddOrderDataDB(context.Background(), ord)
		s.GetCache().Set(ord.Order_uid, ord, 15*time.Minute)

		fmt.Printf("seq = %d [redelivered = %v] mes= %s \n", msg.Sequence, msg.Redelivered, msg.Data)

		time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)

		msg.Ack()

	}, stan.DurableName("i-will-remember"), stan.MaxInflight(100), stan.SetManualAckMode())

	if err != nil {
		log.Fatalln(err)
	}
}

func startNats(s repository.Repository, conf *config.Config) {
	if err := runNats(s, conf); err != nil {
		log.Fatalln(err)
	}
}

func runNats(s repository.Repository, conf *config.Config) error {
	conn, err := stan.Connect(
		conf.Nats.ClusterId,
		conf.Nats.ClientId,
	)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("Connect NATS-Streaming")
	defer logCloser(conn)

	done := make(chan struct{})
	time.Sleep(time.Duration(rand.Intn(4000)) * time.Millisecond)
	SubscriberNats(s, conn)
	<-done

	return nil
}
