package utils

import (
	"testing"
	"strings"
	"net/http/httptest"
	"net/http"
	"log"
	"strconv"

	"github.com/YaNeAndrey/ya-metrics/internal/agent/config"
	"github.com/YaNeAndrey/ya-metrics/internal/storage"
	"github.com/YaNeAndrey/ya-metrics/internal/server/handlers"
	"github.com/YaNeAndrey/ya-metrics/internal/constants"

	"github.com/stretchr/testify/assert"
	"github.com/go-chi/chi/v5"
)

func Test_sendAllMetricsUpdates(t *testing.T) {
	type args struct {
		ms *storage.MemStorage
		c  *config.Config
	}
	tests := []struct {
		name string
		args args
	}{
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sendAllMetricsUpdates(tt.args.ms, tt.args.c)
		})
	}
}

func Test_sendOneMetricUpdate(t *testing.T) {
	type args struct {
		c         *config.Config
		metrType  string
		metrName  string
		metrValue string
	}
	tests := []struct {
		name    string
		args    args
		ms 		*storage.MemStorage
		wantErr bool
	}{
		{
			name: "First test. Update gauge metric",
			args: args {
				c: config.NewConfig(),
				metrType: constants.GaugeMetricType,
				metrName: "gaugeMetric",
				metrValue: "333.3",
			},
			ms: storage.NewMemStorage(),
			wantErr: false,
		},

		{
			name: "Second test. Update counter metric",
			args: args {
				c: config.NewConfig(),
				metrType: constants.CounterMetricType,
				metrName: "counterMetric",
				metrValue: "111",
			},
			ms: storage.NewMemStorage(),
			wantErr: false,
		},
		
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := chi.NewRouter()
			r.Route("/update", func(r chi.Router) {
				r.Post("/{type}/{name}/{value}", func(rw http.ResponseWriter, r *http.Request) {
					handlers.HandlePostUpdateMetricValue(rw,r,tt.ms)
				})
			})

			server := httptest.NewServer(r)
			defer server.Close()
		
			host := strings.Split(server.URL,"/")[2]
			hostBufSlice := strings.Split(host,":")

			hostname := hostBufSlice[0]
			port,_ := strconv.Atoi(hostBufSlice[1])
			
			tt.args.c.SetSrvAddr(hostname)
			tt.args.c.SetSrvPort(port)

			err:= sendOneMetricUpdate(tt.args.c, tt.args.metrType, tt.args.metrName, tt.args.metrValue)
			if err != nil{
				log.Fatal(err)
			}

			switch tt.args.metrType {
				case constants.GaugeMetricType: 
					realValue := strconv.FormatFloat(tt.ms.ListAllGaugeMetrics()[tt.args.metrName], 'f', -1, 64)
					assert.Equal(t,tt.args.metrValue,realValue)

				case constants.CounterMetricType:
					realValue := strconv.FormatInt(tt.ms.ListAllCounterMetrics()[tt.args.metrName], 10)
					log.Println(realValue)
					assert.Equal(t,tt.args.metrValue,realValue)
				}
		})
	}
}

func Test_collectNewMetrics(t *testing.T) {
	type args struct {
		ms *storage.MemStorage
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "First test. Collect all metrics",
			args: args{
				ms: storage.NewMemStorage(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			collectNewMetrics(tt.args.ms)
			assert.Equal(t,28,len(tt.args.ms.ListAllGaugeMetrics()))
			assert.Equal(t,1,len(tt.args.ms.ListAllCounterMetrics()))
		})
	}
}
