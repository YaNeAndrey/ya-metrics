package main

import (
	"fmt"
	"net/http"

	"github.com/YaNeAndrey/ya-metrics/internal/server/router"
	"github.com/YaNeAndrey/ya-metrics/internal/storage"
)

func main() {
	conf := parseFlags()
	ms := storage.NewMemStorage()
	r := router.InitRouter(ms)

	err := http.ListenAndServe(fmt.Sprintf("%s:%d", conf.SrvAddr(), conf.SrvPort()), r)
	if err != nil {
		panic(err)
	}
}
