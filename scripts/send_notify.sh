#!/usr/bin/env bash
grpcurl --plaintext -d "{\"UnitName\": \"$1\"}" localhost:9092 gocron_server.Scheduler/RunJob | jq
