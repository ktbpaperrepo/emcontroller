## What is this program for?
In our experiments, we need to simulate the practical multi-cloud environment, and there are network delays among practical multiple clouds, so we make this program to use TC (Traffic Control) to add network delay among clouds.

### How to generate delays among clouds?
1. Configure the delay values in the function `TestGenCloudsDelay`.
2. At the root directory of this project, run:
```
go test <Project Root Directory>/auto-schedule/experiments/gen-net-delay/ -v -count=1 -run TestGenCloudsDelay
```
For example, if the `<Project Root Directory>` is `/mnt/c/mine/code/gocode/src/emcontroller`, we should execute:
```
go test /mnt/c/mine/code/gocode/src/emcontroller/auto-schedule/experiments/gen-net-delay/ -v -count=1 -run TestGenCloudsDelay
```

### How to clear all delays among clouds?
At the root directory of this project, run:
```
go test <Project Root Directory>/auto-schedule/experiments/gen-net-delay/ -v -count=1 -run TestClearAllDelay
```
For example, if the `<Project Root Directory>` is `/mnt/c/mine/code/gocode/src/emcontroller`, we should execute:
```
go test /mnt/c/mine/code/gocode/src/emcontroller/auto-schedule/experiments/gen-net-delay/ -v -count=1 -run TestClearAllDelay
```
