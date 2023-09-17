import requests
import datetime
import json
import data_types
from types import SimpleNamespace

# endpoint of multi-cloud manager
MCM_END_POINT = "172.27.15.31:20000"
# MCM_END_POINT = "localhost:20000"


def resp_code_successful(code: int) -> bool:
    return code >= 200 and code < 300


def get_all_apps():
    url = "http://" + MCM_END_POINT + "/application"
    headers = {
        'Accept': 'application/json',
    }
    response = requests.get(url, headers=headers, timeout=10)

    if not resp_code_successful(response.status_code):
        raise Exception(
            "URL {}, Unexcepted status code: {}, response body: {}".format(
                url, response.status_code, response.text))

    return json.loads(response.text,
                      object_hook=lambda d: SimpleNamespace(**d))


def call_app(app: data_types.AppInfo) -> data_types.ResultData:
    # endpoint of this app
    app_ep = "{}:{}".format(app.nodePortIP[0], app.nodePort[0])

    url = "http://" + app_ep + "/experiment"
    time_before = datetime.datetime.now()
    response = requests.get(url)
    time_after = datetime.datetime.now()
    if not resp_code_successful(response.status_code):
        raise Exception(
            "URL {}, Unexcepted status code: {}, response body: {}".format(
                url, response.status_code, response.text))
    durations = (time_after - time_before).total_seconds() * 1000  # unit: ms
    return data_types.ResultData(
        app_name=app.appName,
        priority=app.priority,
        resp_time=durations,
        resp_time_in_clouds=float(response.text),
        pri_wei_resp_time=durations * float(app.priority),
        pri_wei_resp_time_in_clouds=float(response.text) * float(app.priority))
