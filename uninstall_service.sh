#!/bin/env bash

systemctl stop emcontroller
systemctl disable emcontroller
rm /usr/lib/systemd/system/emcontroller.service
systemctl daemon-reload
systemctl reset-failed