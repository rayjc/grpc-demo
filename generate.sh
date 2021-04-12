#!/bin/bash

protoc --go_out=plugins=grpc:. --go_opt=paths=source_relative greet/greetpb/greet.proto
protoc --go_out=plugins=grpc:. --go_opt=paths=source_relative calculator/calculatorpb/calculator.proto