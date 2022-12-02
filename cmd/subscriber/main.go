package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/hibiken/asynq"
)

func initiateAsynq() *asynq.Server {
	redisClient := asynq.RedisClientOpt{
		Addr: "0.0.0.0:56379",
	}

	srv := asynq.NewServer(
		redisClient,
		asynq.Config{
			// Specify how many concurrent workers to use
			Concurrency: 10,
			// Optionally specify multiple queues with different priority.
			Queues: map[string]int{"critical": 6, "default": 3, "low": 1},
			// See the godoc for other configuration options
		},
	)

	return srv
}

func setupAsynqHandler() asynq.Handler {
	mux := asynq.NewServeMux()

	mux.HandleFunc("simple-handler", func(ctx context.Context, t *asynq.Task) error {
		log.Printf("simple-handler | payload: %s\n", string(t.Payload()))
		return nil
	})

	mux.HandleFunc("skip-retry", func(ctx context.Context, t *asynq.Task) error {
		log.Printf("skip-retry | payload: %s\n", string(t.Payload()))
		return fmt.Errorf("skip retry %v", asynq.SkipRetry)
	})

	retry := true
	mux.HandleFunc("retry", func(ctx context.Context, t *asynq.Task) error {
		log.Printf("retry | payload: %s retry:%v\n", string(t.Payload()), retry)
		if retry {
			retry = false
			return errors.New("test retry")
		} else {
			retry = true
			return nil
		}
	})

	mux.HandleFunc("handler-pause", func(ctx context.Context, t *asynq.Task) error {
		log.Printf("handler-pause | payload: %s\n", string(t.Payload()))
		time.Sleep(2 * time.Second)
		return nil
	})

	return mux
}

func main() {
	asynqHandler := setupAsynqHandler()
	server := initiateAsynq()
	if err := server.Run(asynqHandler); err != nil {
		log.Fatal("could not run asynq worker")
	}
}
