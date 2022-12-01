package main

import (
	"context"
	"log"

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

	mux.HandleFunc("test", func(c context.Context, t *asynq.Task) error {
		log.Println(string(t.Payload()))
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
