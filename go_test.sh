#!/bin/env bash

# Set the errexit option to exit on error
set -e

CURRENT_DIR=$(cd $(dirname $0); pwd)

# We use "-short" to skip some test functions.
go test ${CURRENT_DIR}/auto-schedule/model/ -count=1 -short
go test ${CURRENT_DIR}/auto-schedule/algorithms/ -count=1 -short
go test ${CURRENT_DIR}/auto-schedule/executors/ -count=1 -short

# the -run parameter of go test reads Regex
# we use the following form to make the code more clear, readable, and maintainable.
funcsToTestInModels="^("
funcsToTestInModels="${funcsToTestInModels}TestLeastRemainPct"
funcsToTestInModels="${funcsToTestInModels}|TestOverflow"
funcsToTestInModels="${funcsToTestInModels}|TestGroupVmsByCloud"
funcsToTestInModels="${funcsToTestInModels}|TestCalcVmAvailVcpu"
funcsToTestInModels="${funcsToTestInModels}|TestCalcVmAvailRamMiB"
funcsToTestInModels="${funcsToTestInModels}|TestCalcVmAvailStorGiB"
funcsToTestInModels="${funcsToTestInModels}|TestCalcVmTotalVcpu"
funcsToTestInModels="${funcsToTestInModels}|TestCalcVmTotalRamMiB"
funcsToTestInModels="${funcsToTestInModels}|TestCalcVmTotalStorGiB"
funcsToTestInModels="${funcsToTestInModels}|TestFindIdxNodeInList"
funcsToTestInModels="${funcsToTestInModels}|TestRemoveNodeFromList"
funcsToTestInModels="${funcsToTestInModels}|TestFindIdxVmInList"
funcsToTestInModels="${funcsToTestInModels}|TestRemoveVmFromList"
funcsToTestInModels="${funcsToTestInModels}|TestGetResOccupiedByPod"
funcsToTestInModels="${funcsToTestInModels})$"

echo "In ${CURRENT_DIR}/models/, the functions to test are ${funcsToTestInModels}."
go test ${CURRENT_DIR}/models/ -count=1 -run "${funcsToTestInModels}"

