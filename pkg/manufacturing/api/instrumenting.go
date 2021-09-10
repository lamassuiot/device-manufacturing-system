package api

import (
	"context"
	"fmt"
	"time"

	"github.com/go-kit/kit/metrics"
)

type instrumentingMiddleware struct {
	requestCount   metrics.Counter
	requestLatency metrics.Histogram
	next           Service
}

func NewInstumentingMiddleware(counter metrics.Counter, latency metrics.Histogram) Middleware {
	return func(next Service) Service {
		return &instrumentingMiddleware{
			requestCount:   counter,
			requestLatency: latency,
			next:           next,
		}
	}
}

func (mw *instrumentingMiddleware) Health(ctx context.Context) bool {
	defer func(begin time.Time) {
		lvs := []string{"method", "Health", "error", "false"}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	return mw.next.Health(ctx)
}

func (mw *instrumentingMiddleware) PostGetCRT(ctx context.Context, keyAlg string, keySize int, c, st, l, o, ou, cn, email, deviceId, caName string) (data []byte, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "PostGetCRT", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	return mw.next.PostGetCRT(ctx, keyAlg, keySize, c, st, l, o, ou, cn, email, deviceId, "")
}

func (mw *instrumentingMiddleware) PostSetConfig(ctx context.Context, authCRT string, CA string) (err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "PostSetConfig", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	return mw.next.PostSetConfig(ctx, authCRT, CA)
}
