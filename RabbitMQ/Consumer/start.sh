#!/bin/sh

./wait-for-it.sh rabbitmq:5672 -t 60    # Wait for RabbitMQ to start

exec go run Consumer.go                      # Run your Go application
