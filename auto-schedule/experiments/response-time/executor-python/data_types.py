from dataclasses import dataclass


# with @dataclass, we do not need to implement the init function of a Class
@dataclass
class ResultData:
    app_name: str
    priority: int
    resp_time: float
    resp_time_in_clouds: float
    pri_wei_resp_time: float  # priority weighted response time
    pri_wei_resp_time_in_clouds: float  # priority weighted response time consumed in clouds


# At first I wanted to perse json directly to PodHost and AppInfo, but then I found that it is complicated, and SimpleNamespace is much simper, so I gave up this.
@dataclass
class PodHost:
    podIP: str
    hostName: str
    hostIP: str


@dataclass
class AppInfo:
    appName: str
    svcName: str
    deployName: str
    clusterIP: str
    nodePortIP: list[str]
    svcPort: list[str]
    nodePort: list[str]
    containerPort: list[str]
    hosts: list[PodHost]
    status: str
    priority: int
    autoScheduled: bool
