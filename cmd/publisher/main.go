package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/hibiken/asynq"
)

const redisAddr = "0.0.0.0:56379"

func main() {
	client := asynq.NewClient(asynq.RedisClientOpt{Addr: redisAddr})
	defer client.Close()

	payload, err := json.Marshal(map[string]interface{}{"payment_id": 10})
	if err != nil {
		log.Fatalf("could not marshal payload: %v", err)
	}

	// run task after 30 second
	{
		task := asynq.NewTask("simple-handler", payload)

		info, err := client.Enqueue(task, asynq.MaxRetry(3), asynq.ProcessIn(3*time.Second), asynq.Queue("critical"), asynq.Timeout(60*time.Second))
		if err != nil {
			log.Fatalf("could not schedule task: %v", err)
		}

		log.Printf("enqueued task: id=%s queue=%s", info.ID, info.Queue)
	}

	// // run task after 60 second and skip retry
	{
		task := asynq.NewTask("skip-retry", payload)

		info, err := client.Enqueue(task, asynq.MaxRetry(3), asynq.ProcessIn(1*time.Second), asynq.Queue("default"))
		if err != nil {
			log.Fatalf("could not schedule task: %v", err)
		}

		log.Printf("enqueued task: id=%s queue=%s", info.ID, info.Queue)
	}

	// distributed task
	{
		for i := 0; i < 100; i++ {
			payload, err := json.Marshal(map[string]interface{}{"test": fmt.Sprintf("test %d", i)})
			if err != nil {
				log.Fatalf("could not marshal payload: %v", err)
			}

			task := asynq.NewTask("handler-pause", payload)

			info, err := client.Enqueue(task, asynq.MaxRetry(3), asynq.Queue("low"))
			if err != nil {
				log.Fatalf("could not schedule task: %v", err)
			}

			log.Printf("enqueued task: id=%s queue=%s", info.ID, info.Queue)
		}
	}

	// test retry
	{
		task := asynq.NewTask("retry", payload)

		info, err := client.Enqueue(task, asynq.MaxRetry(3), asynq.ProcessIn(30*time.Second), asynq.Queue("default"))
		if err != nil {
			log.Fatalf("could not schedule task: %v", err)
		}

		log.Printf("enqueued task: id=%s queue=%s", info.ID, info.Queue)
	}

	// // test no handler
	{
		task := asynq.NewTask("no-handler", payload)

		info, err := client.Enqueue(task, asynq.MaxRetry(0), asynq.ProcessIn(1*time.Second), asynq.Queue("default"))
		if err != nil {
			log.Fatalf("could not schedule task: %v", err)
		}

		log.Printf("enqueued task: id=%s queue=%s", info.ID, info.Queue)
	}
}
