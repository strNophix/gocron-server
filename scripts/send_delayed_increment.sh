#!/usr/bin/env bash
NEW_EPOCH=$(($(date +%s) + 5))
grpcurl --plaintext -d "{\"UnitName\": \"incr\", \"RunAt\": $NEW_EPOCH}" localhost:9092 gocron_server.Scheduler/RunJob | jq
