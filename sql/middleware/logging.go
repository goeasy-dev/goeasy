package middleware

import (
	"goeasy.dev/observability/log"
)

func LogMiddleware(next SQLRequestFunc) SQLRequestFunc {
	return func(request SQLRequest) error {
		log.Debug(request.context, request.query)
		return next(request)
	}
}
