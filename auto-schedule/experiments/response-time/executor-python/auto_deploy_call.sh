#!/bin/env bash

# Set the errexit option to exit on error
set -e

CURRENT_DIR=$(cd $(dirname $0); pwd)

PY_CMD_FILE="python3.11" # In the caller.py, there are multiprocessing and concurrent, maybe only Python 3.11 can do it. I am not sure, but I use python3.11.

MCM_EP="172.27.15.31:20000" # endpoint of multi-cloud manager

declare -i REPEAT_COUNT=30 # use declare to define an integer
ALGO_NAMES=("BERand" "Amaga" "Ampga" "Diktyoga" "Mcssga")
#declare -i REPEAT_COUNT=1 # use declare to define an integer
#ALGO_NAMES=("BERand")
declare -i REQ_COUNT_PER_APP=10

DATA_PATH="${CURRENT_DIR}/data"
JSON_FILE_NAME="request_applications.json"
CALL_PY_FILE="caller.py"
DEL_APPS_PY_FILE="deleter.py"

# a function to print log with time
function print_log() {
  local log_content="$1"
  local log_time=$(date +"%Y-%m-%d %H:%M:%S")
  local caller_info="$(basename ${BASH_SOURCE}):${BASH_LINENO}"
  echo "${log_time} [${caller_info}] ${log_content}"
}

# traverse every repeat
for ((i=1;i<=REPEAT_COUNT;i++))
do
  repeat_path="${DATA_PATH}/repeat${i}"
  json_file_path="${repeat_path}/${JSON_FILE_NAME}"
  # traverse every algorithm
  for algo_name in "${ALGO_NAMES[@]}"
  do
    algo_path="${repeat_path}/${algo_name}"
    print_log "algorithm path is: ${algo_path}"

    curl_cmd="curl -i -X POST -H Content-Type:application/json -H Mcm-Scheduling-Algorithm:${algo_name} -H Expected-Time-One-Cpu:42.629 -d @${json_file_path} http://${MCM_EP}/doNewAppGroup"

     # execute the curl command to deploy applications
    print_log "Execute command: ${curl_cmd}"
    ${curl_cmd}

    # Then, we call the applications

    # request one time to avoid some startup things.
    sleep 10s
    call_cmd="${PY_CMD_FILE} -u ${CURRENT_DIR}/${CALL_PY_FILE}"
    print_log "Execute command: ${call_cmd}"
    ${call_cmd}

    sleep 1s
#    rm "${DATA_PATH}"/*.csv
    rm_csv_cmd="rm ${DATA_PATH}/*.csv"
    print_log "Execute command: ${rm_csv_cmd}"
    ${rm_csv_cmd}

    # request several times to collect data
    for ((j=0; j<REQ_COUNT_PER_APP; j++))
    do
      sleep 1s
      print_log "Execute command: ${call_cmd}"
      ${call_cmd}
    done

    # move the csv files to the repeat and algorithm folder
    sleep 1s
    mv_csv_cmd="mv ${DATA_PATH}/*.csv ${algo_path}"
    print_log "Execute command: ${mv_csv_cmd}"
    ${mv_csv_cmd}

    # delete the deployed applications via multi-cloud manager
    sleep 1s
    del_apps_cmd="${PY_CMD_FILE} -u ${CURRENT_DIR}/${DEL_APPS_PY_FILE}"
    print_log "Execute command: ${del_apps_cmd}"
    ${del_apps_cmd}

  done
done