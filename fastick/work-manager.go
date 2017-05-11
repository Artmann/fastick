// Copyright © 2017 Christoffer Artmann <christoffer@artmann.co>
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package fastick

import (
	"github.com/streadway/amqp"
	"encoding/json"
	"log"
	"fmt"
)

type WorkConfig struct {
	WorkerCount   int

	QueueHost     string
	QueueUsername string
	QueuePassword string
	QueuePort     int
}

type WorkManager struct {
	config WorkConfig
}

func NewWorkManager(config WorkConfig) WorkManager {
	return WorkManager{config: config }
}

func (self *WorkManager) Run() {
	tasks := make(chan Endpoint, self.config.WorkerCount * 5)
	defer close(tasks)

	for w := 1; w <= self.config.WorkerCount; w++ {
		go worker(tasks)
	}
	log.Printf("%d workers started\n", self.config.WorkerCount)

	connectionInfo := fmt.Sprintf("amqp://%s:%s@%s:%d/", self.config.QueueUsername, self.config.QueuePassword, self.config.QueueHost, self.config.QueuePort)
	conn, err := amqp.Dial(connectionInfo)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	taskQueue, err := ch.QueueDeclare(
		"tasks", // name
		true, // durable
		false, // delete when usused
		false, // exclusive
		false, // no-wait
		nil, // arguments
	)
	failOnError(err, "Failed to declare a queue")

	messages, err := ch.Consume(
		taskQueue.Name, // queue
		"", // consumer
		true, // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil, // args
	)

	for message := range messages {
		var endpoint Endpoint
		err := json.Unmarshal(message.Body, &endpoint)
		if err != nil {
			log.Println(err)
		}
		tasks <- endpoint
	}
}