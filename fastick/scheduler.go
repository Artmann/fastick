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
	"log"
	"fmt"
	"database/sql"
	_ "github.com/lib/pq"
	"time"
	"github.com/streadway/amqp"
	"encoding/json"
)

type SchedulerConfig struct {
	Interval         int

	DatabaseHost     string
	DatabaseUsername string
	DatabasePassword string
	DatabaseName     string

	QueueHost        string
	QueueUsername    string
	QueuePassword    string
	QueuePort        int
}

func NewScheduler(config SchedulerConfig) Scheduler {
	return Scheduler{config: config }
}

type Scheduler struct {
	config SchedulerConfig
}

func (self *Scheduler) Run() {
	database := self.initDatabase()
	defer database.Close()

	log.Println("Scheduler running...")
	for {
		rows, err := database.Query("SELECT e.id as ID, e.method AS Method, e.path as Path, a.base_url as BaseURL FROM endpoints e INNER JOIN apps a ON e.app_id = a.id;")
		if err != nil {
			log.Println(err)
		}

		numberOfRows := 0
		for rows.Next() {
			var endpoint Endpoint
			err = rows.Scan(&endpoint.ID, &endpoint.Method, &endpoint.Path, &endpoint.BaseURL)
			if err != nil {
				log.Println(err)
			}

			err := self.addToQueue(endpoint)
			if err != nil {
				log.Println(err)
			}
			numberOfRows++
		}

		log.Printf("Scheduled %d endpoints\n", numberOfRows)
		time.Sleep(time.Duration(self.config.Interval) * time.Second)
	}
}

func (self *Scheduler) addToQueue(endpoint Endpoint) error {
	connectionInfo := fmt.Sprintf("amqp://%s:%s@%s:%d/", self.config.QueueUsername, self.config.QueuePassword, self.config.QueueHost, self.config.QueuePort)
	conn, err := amqp.Dial(connectionInfo)
	if err != nil {
		return err
	}

	ch, err := conn.Channel()
	if err != nil {
		return err
	}

	queue, err := ch.QueueDeclare(
		"tasks", // name
		true, // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil, // arguments
	)
	if err != nil {
		return err
	}

	body, err := json.Marshal(endpoint)
	if err != nil {
		return err
	}

	err = ch.Publish(
		"", // exchange
		queue.Name, // routing key
		false, // mandatory
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         []byte(body),
		})
	if err != nil {
		return err
	}

	return nil
}

func (self *Scheduler) initDatabase() *sql.DB {
	connectionInfo := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable",
		self.config.DatabaseHost, self.config.DatabaseUsername, self.config.DatabasePassword, self.config.DatabaseName)

	db, err := sql.Open("postgres", connectionInfo)
	failOnError(err, "Could not connect to database.")
	return db
}

