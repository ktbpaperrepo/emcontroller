#!/bin/env bash

CURRENT_DIR=$(cd $(dirname $0); pwd)
LOG_DIR=$(cd ${CURRENT_DIR}/..; pwd)

SVC_NAME="emcontroller"

# give execution permission
chmod +x ${CURRENT_DIR}/${SVC_NAME}
chmod +x ${CURRENT_DIR}/run.sh

# render place holders. We should use | instead of / as the delimiter, because the the path contains /
sed -i "s|\[\[PROJECTPATH\]\]|${CURRENT_DIR}|g" ${CURRENT_DIR}/${SVC_NAME}.service
sed -i "s|\[\[LOGPATH\]\]|${LOG_DIR}|g" ${CURRENT_DIR}/${SVC_NAME}.service

# install the service
#cp -f ${CURRENT_DIR}/emcontroller.service /usr/lib/systemd/system/emcontroller.service
cp -f ${CURRENT_DIR}/${SVC_NAME}.service /etc/systemd/system/${SVC_NAME}.service
mkdir -p /etc/systemd/system/${SVC_NAME}.service.d
cp -f ${CURRENT_DIR}/${SVC_NAME}-service.conf /etc/systemd/system/${SVC_NAME}.service.d/${SVC_NAME}-service.conf
systemctl daemon-reload
systemctl enable ${SVC_NAME}
systemctl start ${SVC_NAME}