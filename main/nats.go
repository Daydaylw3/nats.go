package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/nats-io/nats.go"
)

func main() {
	nc, err := nats.Connect("nats://127.0.0.1:4231",
		nats.UserInfo("dayday", "daydaylw333"),
		nats.Timeout(time.Hour*1),
		nats.PingInterval(time.Minute*20),
		nats.SetCustomDialer(&net.Dialer{
			Timeout:   time.Hour * 1,
			KeepAlive: time.Minute * 10,
		}))
	if err != nil {
		log.Printf("%v", err)
		return
	}
	http.HandleFunc("/sayhi", func(w http.ResponseWriter, r *http.Request) {
		err := nc.Publish("foo", []byte("hello world"))
		if err != nil {
			http.Error(w, fmt.Sprintf("%v", err), http.StatusInternalServerError)
		}
	})
	//http.HandleFunc("/help", func(w http.ResponseWriter, r *http.Request) {
	//	msg, err := nc.Request("help", []byte("help me"), time.Minute*3)
	//	if err != nil {
	//		http.Error(w, fmt.Sprintf("%v", err), http.StatusInternalServerError)
	//		return
	//	}
	//	_, _ = w.Write(msg.Data)
	//})
	go func() {
		ch := make(chan *nats.Msg, 64)
		_, err := nc.ChanSubscribe("foo", ch)
		if err != nil {
			log.Printf("%v", err)
			return
		}
		for msg := range ch {
			log.Printf("foo receive msg: %s", msg.Data)
		}
	}()
	//ch := make(chan struct{})
	//go func() {
	//	nc.Subscribe("help", func(m *nats.Msg) {
	//		log.Printf("receive fxxking msg: %s", m.Data)
	//		nc.Publish(m.Reply, []byte("fxxk you man"))
	//	})
	//	<-ch
	//}()
	if err := http.ListenAndServe("localhost:8080", nil); err != nil {
		log.Printf("err: %v", err)
	}
}

func main1() {
	nc, err := nats.Connect("nats://127.0.0.1:4231",
		nats.UserInfo("dayday", "daydaylw333"),
		nats.Timeout(time.Hour*1))
	if err != nil {
		log.Printf("%v", err)
		return
	}
	//err = nc.Publish("foo", []byte("Hello World"))
	msg, err := nc.Request("help", []byte("help me"), 10*time.Second)
	if err != nil {
		log.Printf("%v", err)
		return
	}
	log.Printf("%s", msg.Data)
}
