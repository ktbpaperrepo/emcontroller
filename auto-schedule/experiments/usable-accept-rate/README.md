## What is this program for?
The experiment about:
1. The probability of getting an unusable solution.
2. Priority-weighted acceptance rate.

## How to use this program?
1. Uncomment the debug code in the function `CreateAutoScheduleApps`, but the `draw evolution chart` code should still be commented.
2. Run multi-cloud manager at `localhost:20000` (it will be in the `debug` mode due to last step). The network state database should be running **either** with *this* multi-cloud manager **or** *another* multi-cloud manager.
3. At the root directory of this project, run:
```
go test <Project Root Directory>/auto-schedule/experiments/usable-accept-rate/ -timeout 99999s -v -count=1 -run TestExecute
```
For example, if the `<Project Root Directory>` is `/mnt/c/mine/code/gocode/src/emcontroller`, we should execute:
```
go test /mnt/c/mine/code/gocode/src/emcontroller/auto-schedule/experiments/usable-accept-rate/ -timeout 99999s -v -count=1 -run TestExecute
```
Then, the data file `usable_acceptance_rate.csv` will be generated in this folder.

Then, we can draw a bar chart to compare the application acceptance rate of every algorithm for each priority. In a computer with GUI, in this folder, run:
```
python -u drawer.py
```
