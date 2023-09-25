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

Then, we can:
- use the file `usable_acceptance_rate.csv`, draw a bar chart to compare the application acceptance rate of every algorithm for each priority. In a computer with GUI, in this folder, run:
```
python -u drawer_acc_rate.py
```

We can also do the above step several times **with different numbers of applications** and rename the generated `.csv` files to `usable_acceptance_rate_<number of applications>.csv`, if there are `70` applications in a request, we should rename the file `usable_acceptance_rate.csv` to `usable_acceptance_rate_70.csv`. Then, we can change the value of the constant `APP_COUNTS` in `drawer_sched_time.py` and `drawer_total_acc_rate.py`, according to the different numbers of applications.
Then, we can:
- draw a bar chart to compare the application acceptance rate of every algorithm with different numbers of applications. In a computer with GUI, in this folder, run:
```
python -u drawer_total_acc_rate.py
```
- draw a bar chart to compare the maximum scheduling time used by every algorithm with different numbers of applications. In a computer with GUI, in this folder, run:
```
python -u drawer_sched_time.py
```