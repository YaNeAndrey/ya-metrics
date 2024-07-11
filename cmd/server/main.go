package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/YaNeAndrey/ya-metrics/internal/server/router"
	"github.com/YaNeAndrey/ya-metrics/internal/server/utils"
	"github.com/YaNeAndrey/ya-metrics/internal/storage"
	"github.com/YaNeAndrey/ya-metrics/internal/storage/storagedb"
	"github.com/YaNeAndrey/ya-metrics/internal/storage/storagejson"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	pb "github.com/YaNeAndrey/ya-metrics/internal/proto"
	myrpc "github.com/YaNeAndrey/ya-metrics/internal/server/grpc"
	_ "net/http/pprof"
)

var buildVersion = "N/A"
var buildDate = "N/A"
var buildCommit = "N/A"

func main() {

	log.Printf("Build version: %s", buildVersion)
	log.Printf("Build date: %s", buildDate)
	log.Printf("Build commit: %s", buildCommit)

	log.SetReportCaller(true)

	conf := parseFlags()
	testMetrics := []storage.Metrics{}

	var st storage.StorageRepo
	var err error
	if conf.DBConnectionString() != "" {
		st, err = storagedb.InitStorageDB(conf.DBConnectionString())
		if err != nil {
			log.Println(err)
		}
	}

	if st == nil {
		st = storage.StorageRepo(storagejson.NewMemStorageJSON(testMetrics))

		err = utils.ReadMetricsFromFile(conf.FileStoragePath(), &st)
		if err != nil {
			log.Println(err.Error())
		}

		if conf.StoreInterval() != 0 {
			go utils.SaveMetricsByTime(conf.FileStoragePath(), conf.StoreInterval(), &st)
		}
		defer utils.SaveAllMetricsToFile(conf.FileStoragePath(), &st)
	}

	log.Printf(conf.String())

	var srv = http.Server{Addr: conf.SrvAddr()}
	srv.Handler = router.InitRouter(*conf, &st)

	idleConnsClosed := make(chan struct{})
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	go func() {
		<-sigint
		if err := srv.Shutdown(context.Background()); err != nil {
			log.Printf("HTTP server Shutdown: %v", err)
		}
		close(idleConnsClosed)
	}()

	if conf.RPCAddr() != "" {
		go func() {
			listen, err := net.Listen("tcp", conf.RPCAddr())
			if err != nil {
				log.Fatal(err)
			}
			s := grpc.NewServer()
			pb.RegisterMetricsServer(s, myrpc.NewGRPCMetricsServer(&st))

			fmt.Println("Сервер gRPC начал работу")
			if err := s.Serve(listen); err != nil {
				log.Fatal(err)
			}
		}()
	}

	if err = srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		panic(err)
	}

	<-idleConnsClosed
	fmt.Println("Server Shutdown gracefully")
}
