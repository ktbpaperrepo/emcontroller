#!/bin/env bash

CURRENT_DIR=$(cd $(dirname $0); pwd)

# We use "-short" to skip some test functions.
go test ${CURRENT_DIR}/auto-schedule/model/ -count=1 -short
go test ${CURRENT_DIR}/auto-schedule/algorithms/ -count=1 -short
go test ${CURRENT_DIR}/auto-schedule/executors/ -count=1 -short

go test ${CURRENT_DIR}/models/ -v -count=1 -run TestLeastRemainPct
go test ${CURRENT_DIR}/models/ -v -count=1 -run TestOverflow
go test ${CURRENT_DIR}/models/ -v -count=1 -run TestCalcVmAvailVcpu
go test ${CURRENT_DIR}/models/ -v -count=1 -run TestCalcVmAvailRamMiB
go test ${CURRENT_DIR}/models/ -v -count=1 -run TestCalcVmAvailStorGiB
go test ${CURRENT_DIR}/models/ -v -count=1 -run TestCalcVmTotalVcpu
go test ${CURRENT_DIR}/models/ -v -count=1 -run TestCalcVmTotalRamMiB
go test ${CURRENT_DIR}/models/ -v -count=1 -run TestCalcVmTotalStorGiB
