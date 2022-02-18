package main

import (
	"context"
	"time"

	"github.com/go-kit/kit/log"
)

type LogMiddleware struct {
	logger log.Logger
	CounterService
}

func DefaultLoggingMiddleware(logger log.Logger) CounterServiceMiddleware {
	return func(next CounterService) CounterService {
		return LogMiddleware{logger, next}
	}
}

func (mw LogMiddleware) Init() (err error) {
	defer func(begin time.Time) {
		mw.logger.Log(
			"method", "Init",
			"error", err,
			"took", time.Since(begin))
	}(time.Now())

	err = mw.CounterService.Init()
	return
}

func (mw LogMiddleware) Terminate() (err error) {
	defer func(begin time.Time) {
		mw.logger.Log(
			"method", "Terminate",
			"error", err,
			"took", time.Since(begin))
	}(time.Now())

	err = mw.CounterService.Terminate()
	return
}

func (mw LogMiddleware) Increment(context context.Context, uuid string) (c Counter, err error) {
	defer func(begin time.Time) {
		mw.logger.Log(
			"method", "Increment",
			"input", uuid,
			"output.uuid", c.UUID,
			"output.name", c.Name,
			"output.error", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	c, err = mw.CounterService.Increment(context, uuid)
	return
}

func (mw LogMiddleware) Create(context context.Context, uuid string, name string) (c Counter, err error) {
	defer func(begin time.Time) {
		mw.logger.Log(
			"method", "Create",
			"input", uuid,
			"output.uuid", c.UUID,
			"output.name", c.Name,
			"error", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	c, err = mw.CounterService.Create(context, uuid, name)
	return
}
