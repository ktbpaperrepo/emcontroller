## What is this program for?
Automatic experiments of response time.

## How to execute the experiments using this program?
### Step: Decide "repeat count", "device count", "app count", "request count per app", "multi-cloud manager endpoint", "expt-app name prefix", "auto-scheduling VM name prefix" and maybe others.

Currently, if we want the experiments automatic, `device count` should be `1`. 

After deciding them, we should set them in the related files including `init.go`, `auto_deploy_call.sh`, `caller.py`, `deleter.py`, `charts_drawer.py`, `http_api.py`.

#### repeat count
We should repeat the experiment multiple times to reduce the impact of random factors.

#### device count, app count, request count per app
`app count` is the number of applications in the deployment/scheduling request. When we use `auto-schedule/experiments/applications-generator` to generate deployment requests, we should set this value.

`device count` and `request count per app`: In every repeat, we can use multiple devices to request the applications to simulate the real production, and each device can access each app several times to reduce the impact of random factors.

### Step: Clear old data
- move all `executor-python/data/repeatX` into the folder `executor-python/data/bak`.
- Delete all `executor-python/data/repeatX`.

### Step: Generate applications deployment request json body and the needed folders
At this folder, run `go run init.go`.

### Step: Deploy applications and call applications for the repeats set by us
At the folder `executor-python`, run `bash auto_deploy_call.sh`.

This step can be executed on a VM.

### Step: Draw charts
On a system wich GUI, at the folder `executor-python`,
- run `python -u charts_drawer.py` to draw cdf charts.
- run `python -u charts_drawer_no_cdf.py` to draw dot charts.