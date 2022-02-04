#!/usr/bin/env bash
grpcurl --plaintext -d "{}" localhost:9092 gocron_server.Scheduler/ListenJobs