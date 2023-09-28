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


# calculate the average value of a given list of ResultData
def calc_rd_avg(data_list: list[ResultData]) -> ResultData:
    app_name = data_list[0].app_name
    priority = data_list[0].priority

    # calculate the average values
    resp_time = 0.0
    resp_time_in_clouds = 0.0
    pri_wei_resp_time = 0.0
    pri_wei_resp_time_in_clouds = 0.0
    for _, one_data in enumerate(data_list):
        resp_time += one_data.resp_time
        resp_time_in_clouds += one_data.resp_time_in_clouds
        pri_wei_resp_time += one_data.pri_wei_resp_time
        pri_wei_resp_time_in_clouds += one_data.pri_wei_resp_time_in_clouds
    resp_time /= len(data_list)
    resp_time_in_clouds /= len(data_list)
    pri_wei_resp_time /= len(data_list)
    pri_wei_resp_time_in_clouds /= len(data_list)

    return ResultData(app_name=app_name,
                      priority=priority,
                      resp_time=resp_time,
                      resp_time_in_clouds=resp_time_in_clouds,
                      pri_wei_resp_time=pri_wei_resp_time,
                      pri_wei_resp_time_in_clouds=pri_wei_resp_time_in_clouds)


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
