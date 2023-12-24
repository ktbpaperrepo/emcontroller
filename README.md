# Multi-Cloud Manager #

## Basic information

This project is developed and tested with `Go version go1.20.2 linux/amd64`.

The project is developed using the Beego framework, with APIs defined in the `routers` and `controllers`, core functionality outlined in the `models` directory, and front-end code located within the `views` and `static` directories. Users are required to configure parameters in the `conf` directory."

### How to build Multi-cloud Manager and clean the built files? ###

* `make` or `make emcontroller` will generate the binary `emcontroller`.
* `make clean` will remove the binary `emcontroller`.

### How do I run and stop Multi-Cloud Manager? ###

* After building the binary `emcontroller`, in the root path of this project, execute `./emcontroller`.
* `Ctrl + C` to stop.

### How do I set Multi-cloud Manager as a service of systemd and delete the service? ###

* After building the binary `emcontroller`, in the root path of this project, execute `bash install_service.sh`.
* Execute `bash uninstall_service.sh` to delete the service.


## Automatic scheduling
Multi-cloud Manager enables the scheduling of applications, as detailed in the paper "_Multi-cloud Containerized Service Scheduling Optimizing Computation and Communication_". This functionality requires information on the Network Round-Trip Time (RTT) between pairs of clouds. To facilitate this, users are required to upload the "network performance test container image" to the container image repository. Multi-cloud Manager employs a periodic task for collecting RTT data.

### How do I make network performance test container image? ###
1. Put the folder `net-perf-container-image` to a VM with Docker installed.
2. On that VM, `cd` into the folder `net-perf-container-image`, and execute `docker build -t mcnettest:latest .`.

The code for automatic scheduling can be found in the `auto-schedule` folder. In particular, the scheduling algorithms used in the "Evaluation section" of the paper are implemented in the following files within the `auto-schedule/algorithms` folder: `mcssga.go`, `for_cmp_amaga.go`, `for_cmp_ampga.go`, `for_cmp_best_effort_rand.go`, and `for_cmp_diktyo_ga.go`.

For the "Dummy Service" discussed in the paper, you can find its code located in the `auto-schedule/experiments/server` folder. Furthermore, the services parameters (e.g., requirements and others) employed in the paper's experiments are generated using the code available in the `auto-schedule/experiments/applications-generator` directory. To access the specific code for two experiments, please navigate to the `auto-schedule/experiments/usable-accept-rate` and `auto-schedule/experiments/response-time` folders. You'll find detailed information provided in the `README.md` file within each respective folder.

## Data of the experiments in paper "_Multi-cloud Containerized Service Scheduling Optimizing Computation and Communication_"
- The data of experiments about Scheduling Time, Usable Solution Rate,and Service Acceptance Rate are the `.csv` files in the folder `auto-schedule/experiments/usable-accept-rate`.
- The data and service groups of experiments about Response Time are in the folder `auto-schedule/experiments/response-time/executor-python/data`.
  - The service groups are `request_applications.json` files.
  - The data are `.csv` files.