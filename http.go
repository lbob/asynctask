package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
)

var workerNum = flag.Int("num", 10, "worker number")
var baseurl = flag.String("baseurl", "", "base url")
var addr = flag.String("addr", ":http", "listen addr:port")

var hub *Scheduler

func main() {
	flag.Parse()

	std := log.New(os.Stdout, "[Info]", log.LstdFlags)
	err := log.New(os.Stderr, "[Scheduler]", log.LstdFlags)

	hub = new(Scheduler).Init(*workerNum, *baseurl, std, err)

	http.HandleFunc("/", page_index)
	http.HandleFunc("/status", page_status)
	http.HandleFunc("/task/add", page_task_add)
	go log.Fatal(http.ListenAndServe(*addr, nil))

	go exitSignal()

	hub.Run()
}

func exitSignal() {
	co := make(chan os.Signal, 1)
	signal.Notify(co, os.Interrupt)
	<-co

	hub.Close()
}

type Result struct {
	Code int
	Data interface{}
}

func page_index(w http.ResponseWriter, r *http.Request) {
	rs := &Result{Code: 0, Data: "ok"}
	json.NewEncoder(w).Encode(rs)
}

func page_task_add(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("action")
	data := r.FormValue("params")

	hub.AddOrder(name, data)

	rs := &Result{Code: 0, Data: "ok"}
	json.NewEncoder(w).Encode(rs)
}

func page_status(w http.ResponseWriter, r *http.Request) {
	t := hub.Status()

	rs := &Result{Code: 0, Data: t}
	json.NewEncoder(w).Encode(rs)
}
