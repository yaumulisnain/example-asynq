package main

import (
	"encoding/json"
	"log"
	"time"

	"github.com/hibiken/asynq"
)

const redisAddr = "0.0.0.0:56379"

func main() {
	client := asynq.NewClient(asynq.RedisClientOpt{Addr: redisAddr})
	defer client.Close()

	payload, err := json.Marshal(map[string]interface{}{"test": "test"})
	if err != nil {
		log.Fatalf("could not marshal payload: %v", err)
	}

	task := asynq.NewTask("test", payload)

	info, err := client.Enqueue(task, asynq.MaxRetry(0), asynq.ProcessIn(3*time.Second), asynq.Queue("critical"))
	if err != nil {
		log.Fatalf("could not schedule task: %v", err)
	}
	log.Printf("enqueued task: id=%s queue=%s", info.ID, info.Queue)
}
