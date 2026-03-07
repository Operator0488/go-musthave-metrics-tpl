package handler

import (
	"context"
	"errors"
	"github.com/go-resty/resty/v2"
	"github.com/jackc/pgx/v5/pgconn"
	"strings"
	"time"
)

func Retry(req *resty.Request, url string) (resp *resty.Response, err error) {
	var timePeriod = []time.Duration{1 * time.Second, 3 * time.Second, 5 * time.Second}

	for _, dur := range timePeriod {
		resp, err = req.Post(url)
		if isRetryHTTP(resp.StatusCode(), err) {
			time.Sleep(dur)
			continue
		}
		break
	}
	return resp, err
}

func isRetryHTTP(status int, err error) bool {
	if err == nil {
		return false
	}

	if errors.Is(err, context.DeadlineExceeded) {
		return true
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return strings.HasPrefix(pgErr.Code, "08")
	}

	switch status {
	case 500, 502, 503, 504:
		return true
	default:
		return false
	}

}
