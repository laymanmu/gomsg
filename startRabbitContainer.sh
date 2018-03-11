#!/usr/bin/env bash

NAME=rabbitmq

docker run -d -p 5672:5672 --rm --hostname $NAME --name $NAME $NAME

