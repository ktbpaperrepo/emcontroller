#!/bin/env bash

CURRENT_DIR=$(cd $(dirname $0); pwd)

go test ${CURRENT_DIR}/auto-schedule/model/ -count=1
go test ${CURRENT_DIR}/auto-schedule/algorithms/ -count=1
go test ${CURRENT_DIR}/auto-schedule/executors/ -count=1
go test ${CURRENT_DIR}/models/ -v -count=1 -run TestLeastRemainPct
go test ${CURRENT_DIR}/models/ -v -count=1 -run TestAllMoreThan