package grpc

import (
	"context"
	"github.com/YaNeAndrey/ya-metrics/internal/constants"
	pb "github.com/YaNeAndrey/ya-metrics/internal/proto"
	"github.com/YaNeAndrey/ya-metrics/internal/storage"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func NewGRPCMetricsServer(st *storage.StorageRepo) *MetricsServer {
	return &MetricsServer{
		storage: st,
	}
}

type MetricsServer struct {
	pb.UnimplementedMetricsServer
	storage *storage.StorageRepo
}

func (ad *MetricsServer) UpdateMetric(ctx context.Context, req *pb.MetricRequest) (*pb.MetricResponse, error) {
	newDelta := req.Metric.Delta
	newValue := req.Metric.Value

	newMetric := storage.Metrics{
		ID: req.Metric.ID,
	}

	if req.Metric.MType == 0 {
		newMetric.MType = constants.GaugeMetricType
		newMetric.Value = &newValue
	} else {
		newMetric.MType = constants.CounterMetricType
		newMetric.Delta = &newDelta
	}

	err := (*ad.storage).UpdateOneMetric(ctx, newMetric, false)
	if err != nil {
		return nil, status.Error(codes.FailedPrecondition, err.Error())
	}

	return nil, nil
}
