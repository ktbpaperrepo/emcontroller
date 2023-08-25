## What is this program for?
This is a server application for experiments about response time.

When this application runs, it will occupy the input values of Storage and Memory. When it receives an API request, it will send a request to every of its dependent applications and wait for the responses, and after that do its own workload, and finally send response.

## How to use this program?
### How to build this application?
1. In this directory, run `go build -o experiment-app`, and then a file named `experiment-app` will be generated. Which is the executable binary file of this application.
2. Run `docker build -t mcexp:latest .`, and then the container image with "RepoTag" `mcexp:latest` will be generated.

### How to run this application?

The parameters are described before the `main()` function in `main.go`.

An example to call multi-cloud manager to run this containerized application is:
```shell
curl -i -X POST -H Content-Type:application/json -d '{"name":"exp-app2","replicas":1,"hostNetwork":false,"nodeName":"testmem","containers":[{"name":"exp-app2","image":"172.27.15.31:5000/mcexp:latest","workDir":"","resources":{"limits":{"memory":"5000Mi","cpu":"2","storage":"10Gi"},"requests":{"memory":"5000Mi","cpu":"0.5","storage":"10Gi"}},"commands":["./experiment-app"],"args":["5000000","2","5000","10","http://exp-app1-service:81/experiment"],"env":null,"mounts":null,"ports":[{"containerPort":3333,"name":"tcp","protocol":"tcp","servicePort":"81","nodePort":"30002"}]}],"priority":0,"autoScheduled":false}' http://172.27.15.31:20000/doNewApplication
```