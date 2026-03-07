package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"net"
	"strconv"
	"time"
)

func getValueFloat(str string) (float64, error) {
	return strconv.ParseFloat(str, 64)
}

func getValueInt(str string) (int64, error) {
	return strconv.ParseInt(str, 10, 64)
}

func retry(ctx context.Context, fn func() error) error {
	var err error
	var timePeriod = []time.Duration{1 * time.Second, 3 * time.Second, 5 * time.Second}

	for _, dur := range timePeriod {
		if err = fn(); isRetryDB(err) {
			select {
			case <-time.After(dur):
			case <-ctx.Done():
				return ctx.Err()
			}
		} else {
			break
		}
	}

	return err
}

func isRetryDB(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
		if pgErr.Code[:2] == "08" {
			fmt.Println("Connection error:", pgErr.Code)
			return true
		}
	}

	var netErr net.Error
	if errors.As(err, &netErr) && netErr.Temporary() {
		return true
	}

	return false
}
