## What is this program for?
The experiment about:
1. The probability of getting an unusable solution.
2. Application acceptance rate.
3. Application acceptance rate of every priority.
4. Scheduling time used by different algorithms.

## How to use this program?
1. Set the value of variables `appCounts` and `repeatCount` in file `executor.go`. Also, set the same application count value to constant `APP_COUNTS` in file `common.py`.
2. Uncomment the debug code in the function `CreateAutoScheduleApps`, but the `draw evolution chart` code should still be commented.
3. Run multi-cloud manager at `localhost:20000` (This is to execute `go run main.go` at the root directory of this project) (it will be in the `debug` mode due to last step). The network state database should be running **either** with *this* multi-cloud manager **or** *another* multi-cloud manager.
4. At this folder, run:
```
go run executor.go
```
Then, the data **files** `usable_acceptance_rate_<appCount>.csv` will be generated in this folder.

We can also compile the code and run it on a VM without Golang installed, which is closer to the real production scenario. To do this, we just need to compile `multi-cloud manager` and this program `usable-accept-rate` after the above _Step 1_. Then, we should `scp` the whole `multi-cloud manager` project and the configuration `/root/.kube/config` (we can simply copy the folder `/root/.kube`) to the VM, and then do 2 and 3 at that VM by executing the **compiled binary files** instead of `go run xxx.go`. If we need to do this experiment for a long time, we need to execute the binary files in the background `nohup ./emcontroller 2>&1 &` and `nohup ./usable-accept-rate 2>&1 &`. Lastly, we need to `scp` the generated `usable_acceptance_rate_<appCount>.csv` files back to this folder from the VM.

Then, we can:
- use the `usable_acceptance_rate_<appCount>.csv` files to draw bar charts to compare the application acceptance rate of every algorithm for each priority, one chart for each application count. In a computer with GUI, in this folder, run:
```
python -u drawer_acc_rate.py
```
- draw a bar chart to compare the application acceptance rate of every algorithm with different numbers of applications. In a computer with GUI, in this folder, run:
```
python -u drawer_total_acc_rate.py
```
- draw a bar chart to compare the maximum scheduling time used by every algorithm with different numbers of applications. In a computer with GUI, in this folder, run:
```
python -u drawer_sched_time.py
```