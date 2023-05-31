#!/bin/env bash

SVC_NAME="emcontroller"

systemctl stop ${SVC_NAME}
systemctl disable ${SVC_NAME}
#rm /usr/lib/systemd/system/emcontroller.service
rm /etc/systemd/system/${SVC_NAME}.service
rm -rf /etc/systemd/system/${SVC_NAME}.service.d
systemctl daemon-reload
systemctl reset-failed