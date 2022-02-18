package main

import (
	"context"
	"fmt"
	"time"

	"github.com/go-kit/kit/metrics"
)

func DefaultInstrumentingMiddleware(
	requestCount metrics.Counter,
	requestLatency metrics.Histogram,
) CounterServiceMiddleware {
	return func(next CounterService) CounterService {
		return InstrumentingMiddleware{requestCount, requestLatency, next}
	}
}

type InstrumentingMiddleware struct {
	requestCount   metrics.Counter
	requestLatency metrics.Histogram
	CounterService
}

func (mw InstrumentingMiddleware) Init() (err error) {
	return mw.CounterService.Init()
}

func (mw InstrumentingMiddleware) Create(context context.Context, uuid string, name string) (c Counter, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "Create", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	c, err = mw.CounterService.Create(context, uuid, name)
	return
}

func (mw InstrumentingMiddleware) Increment(context context.Context, uuid string) (c Counter, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "Increment", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	c, err = mw.CounterService.Increment(context, uuid)
	return
}

func (mw InstrumentingMiddleware) Terminate() error {
	return mw.CounterService.Terminate()
}
