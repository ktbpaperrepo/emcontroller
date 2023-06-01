#!/bin/env bash

# example:
# bash client.sh "192.168.100.136" "NOKIA7" "192.168.100.136" "3306" "multicloud" 'AAUproxmox1234!@#' "multi_cloud" "NOKIA8"

# One simple way to make your script exit on any error is to use the `set -e` option. This will cause your script to immediately exit if any command returns a non-zero exit code. The default value is not sure. `set +e` can cancel it.
set -e

t_cloud_ip="$1"
t_cloud_name="$2"

mysql_ip="$3"
mysql_port="$4"
mysql_user="$5"
mysql_passwd="$6"
mysql_db_name="$7"
this_cloud_name="$8"

# measure the rtt between this cloud an the target cloud
last_line_ping=$(ping -c 10 -w 30 "${t_cloud_ip}" | tail -1)

# if reachable, the last line is like:
# round-trip min/avg/max/stddev = 0.364/0.458/0.739/0.101 ms
# contains '/' and '='
# if unreachable, the last line is like:
# 10 packets transmitted, 0 packets received, 100% packet loss
# does not contain  '/' and '='
rtt_ms=""
if [[ "${last_line_ping}" == *=* && "${last_line_ping}" == */* ]]
then
  rtt_ms=$(echo "${last_line_ping}" | awk '{print $4}' | cut -d '/' -f 2)
  echo "RTT from ${this_cloud_name} to ${t_cloud_name} is ${rtt_ms} ms."
else
  rtt_ms="250000" # very big value means unreachable
  echo "Unreachable from ${this_cloud_name} to ${t_cloud_name}, so we set the RTT as ${rtt_ms} ms."
fi

# write the rtt in to the database
mysql -u "${mysql_user}" --port "${mysql_port}" -h "${mysql_ip}" --password="${mysql_passwd}" -e "update ${mysql_db_name}.${this_cloud_name} set rtt_ms=${rtt_ms} where target_cloud_name='${t_cloud_name}'"
echo "Write the rtt value in the database successfully."