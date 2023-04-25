#!/bin/env bash

CURRENT_DIR=$(cd $(dirname $0); pwd)

# give execution permission
chmod +x ${CURRENT_DIR}/emcontroller
chmod +x ${CURRENT_DIR}/run.sh

# render place holders. We should use | instead of / as the delimiter, because the the path contains /
sed -i "s|\[\[PROJECTPATH\]\]|${CURRENT_DIR}|g" ${CURRENT_DIR}/emcontroller.service

# install the service
cp -f ${CURRENT_DIR}/emcontroller.service /usr/lib/systemd/system/emcontroller.service
systemctl daemon-reload
systemctl enable emcontroller
systemctl start emcontroller