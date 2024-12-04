package metrics

import "github.com/prometheus/client_golang/prometheus"

const (
	orderLabel   = "order"
	handlerLabel = "handler"
	codeLabel    = "code"
)

var (
	ordersGiven = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "orders_created",
			Help: "Количество созданных заказов",
		},
		[]string{orderLabel},
	)

	okResponseTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "ok_response_total",
			Help: "Суммарное количество успешных ответов",
		},
		[]string{handlerLabel},
	)

	errResponseTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "err_response_total",
			Help: "Суммарное количество плохих ответов",
		},
		[]string{handlerLabel, codeLabel},
	)
)

func init() {
	prometheus.MustRegister(ordersGiven, okResponseTotal, errResponseTotal)
}

func IncOkResponseTotal(handler string) {
	okResponseTotal.With(prometheus.Labels{
		handlerLabel: handler,
	}).Inc()
}

func IncErrResponseTotal(handler, code string) {
	errResponseTotal.With(prometheus.Labels{
		handlerLabel: handler,
		codeLabel:    code,
	}).Inc()
}

func AddOrdersGiven(k int) {
	ordersGiven.With(prometheus.Labels{}).Add(float64(k))
}
