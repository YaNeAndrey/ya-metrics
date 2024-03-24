package constants

import "errors"

const (
	GaugeMetricType   = "gauge"
	CounterMetricType = "counter"
)

var ErrIncorrectEndpointFormat = errors.New("need address in a form host:port")
var ErrIncorrectPortNumber = errors.New("port accepts values from the range [1:65535]")
var ErrIncorrectPollInterval = errors.New("pollInterval must be greater than 0")
var ErrIncorrectReportInterval = errors.New("reportInterval must be greater than 0")
var ErrIncorrectStoreInterval = errors.New("StoreInterval must be greater then -1")
var ErrGaugeValue = errors.New("no value specified for gauge metric")
var ErrGaugeDelta = errors.New("delta must be nil for counter metric")
var ErrCounterDelta = errors.New("no delta specified for counter metric")
var ErrCounterValue = errors.New("value must be nil for counter metric")
var ErrIncorectMetricType = errors.New("incorrect metric type")
var ErrMetricTypeDoNotMatch = errors.New("metric types do not match")
var ErrIncorrectRateLimit = errors.New("rate limit must be greater than 0")
