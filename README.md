# README #

Edge and Multi-Cloud Manager.

### How to build Multi-cloud Manager and clean the built files? ###

* `make` or `make emcontroller` will generate the binary `emcontroller`.
* `make clean` will remove the binary `emcontroller`.

### How do I run and stop Multi-Cloud Manager? ###

* After building the binary `emcontroller`, in the root path of this project, execute `./emcontroller`.
* `Ctrl + C` to stop.

### How do I set Multi-cloud Manager as a service of systemd and delete the service? ###

* After building the binary `emcontroller`, in the root path of this project, execute `bash install_service.sh`.
* Execute `bash uninstall_service.sh` to delete the service.

### How do I make network performance test container image? ###
1. Put the folder `net-perf-container-image` to a VM with Docker installed.
2. On that VM, `cd` into the folder `net-perf-container-image`, and execute `docker build -t mcnettest:latest .`.