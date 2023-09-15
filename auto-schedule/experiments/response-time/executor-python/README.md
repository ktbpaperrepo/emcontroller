## What is this program for?
This is a tool to call the experiment applications, record the response time related data, and draw charts with that data. 

This is the main executor of the experiments.

## How to execute the experiments using this program?
### Step: Decide "repeat count", "device count", "app count", "request count per app"

After deciding them, we should set the variables `REPEAT_COUNT`, `DEVICE_COUNT`, `APP_COUNT`, `REQ_COUNT_PER_APP` in the file `charts_drawer.py`.

#### repeat count
We should repeat the experiment multiple times to reduce the impact of random factors.

If we need to repeat the experiment 5 times, we should create 5 folders named `repeat1`, `repeat2`, ..., `repeat5` in the folder `data`.

In each of the `repeatX` folders, we should create:
- folders with the name of the algorithms that we need to compare in our experiments. 
- a file named `request_applications.json` used to store the applications deployment request json body. (This is because in each repeat, we should use the same applications to compare different algorithms.)

#### device count, app count, request count per app
`app count` is the number of applications in the deployment/scheduling request. When we use `auto-schedule/experiments/applications-generator` to generate deployment requests, we should set this value.

`device count` and `request count per app`: In every repeat, we can use multiple devices to request the applications to simulate the real production, and each devices can access each app several times to reduce the impact of random factors.

### Step: Generate applications deployment request json body
Use `auto-schedule/experiments/applications-generator` to generate deployment requests. If `repeat count` is 3, we should generate 3 json bodies and put them in `data/repeatX/request_applications.json`.

### Step: Deploy applications
Send requests to multi-cloud manager to schedule and deploy applications, using the json body in `data/repeatX/request_applications.json`, and use `Mcm-Scheduling-Algorithm` HTTP header to set the name of algorithm.

### Step: Call applications
After applications are deployed and become running, we run `python3.11 -u caller.py` `request count per app` times to access the applications and generate `request count per app`csv files in the folder `data`. 

Then we move all the generated `csv` files to the corresponding folders `data/repeatX/{algorithm name}`.

### Step: Draw charts
On a system wich GUI, run `python -u charts_drawer.py` to draw charts.