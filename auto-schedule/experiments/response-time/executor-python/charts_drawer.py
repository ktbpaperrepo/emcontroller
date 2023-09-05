import glob
import csv
import numpy as np
import matplotlib.pyplot as plt

import data_types
import csv_operation

REPEAT_COUNT = 2  # We repeat the experiments for REPEAT_COUNT times.
DEVICE_COUNT = 1  # we use DEVICE_COUNT devices to send requests.
APP_COUNT = 60  # In every repeat, we deploy APP_COUNT applications.
REQ_COUNT_PER_APP = 5  # In every repeat, on every device, we access every application REQ_COUNT_PER_APP times.
ALGO_NAMES = [
    "BERand", "Ampga", "Mcssga"
]  # the names of all algorithms to be evaluated in this experiment

ATTR_TO_METRIC: dict[str, str] = {
    "resp_time":
    "response time (ms)",
    "resp_time_in_clouds":
    "response time in clouds (ms)",
    "pri_wei_resp_time":
    "priority-weighted response time (ms)",
    "pri_wei_resp_time_in_clouds":
    "priority-weighted response time in clouds (ms)",
}


# calculate the values used to complement the rejected applications.
def calc_values_for_rejected(
    existing_data: dict[str, list[data_types.ResultData]]
) -> data_types.ResultData:

    # initialization
    data_for_rej = data_types.ResultData(app_name="rej",
                                         priority=1,
                                         resp_time=0,
                                         resp_time_in_clouds=0,
                                         pri_wei_resp_time=0,
                                         pri_wei_resp_time_in_clouds=0)

    # set the values in data_for_rej as the largest value of every metric among all algorithms
    for algo_name, algo_data in existing_data.items():
        for _, one_data in enumerate(algo_data):
            if one_data.resp_time > data_for_rej.resp_time:
                data_for_rej.resp_time = one_data.resp_time

            if one_data.resp_time_in_clouds > data_for_rej.resp_time_in_clouds:
                data_for_rej.resp_time_in_clouds = one_data.resp_time_in_clouds

            if one_data.pri_wei_resp_time > data_for_rej.pri_wei_resp_time:
                data_for_rej.pri_wei_resp_time = one_data.pri_wei_resp_time

            if one_data.pri_wei_resp_time_in_clouds > data_for_rej.pri_wei_resp_time_in_clouds:
                data_for_rej.pri_wei_resp_time_in_clouds = one_data.pri_wei_resp_time_in_clouds

            # TODO: Perhaps, for the priority-weighted metrics, we may calculate them with their own priorities instead of directly using the highest priority-weighted values, to show that we only reject the low-proirity applicatioins. However, this may not be necessary, because we already have another metric: priority-weighted acceptance rate.

    # to make a difference between the rejected applicatioins and accepted ones, we enlarge the values in data_for_rej 1.1 times.
    enlarge_rate = 1.1
    data_for_rej.resp_time *= enlarge_rate
    data_for_rej.resp_time_in_clouds *= enlarge_rate
    data_for_rej.pri_wei_resp_time *= enlarge_rate
    data_for_rej.pri_wei_resp_time_in_clouds *= enlarge_rate

    return data_for_rej


# get one metric for making CDF charts
def make_cdf_data(all_data: dict[str, list[data_types.ResultData]],
                  metric_attr_name: str) -> dict[str, list[float]]:

    metric_data: dict[str, list[float]] = dict()

    for algo_name, algo_data in all_data.items():
        metric_data[algo_name] = []
        for _, one_data in enumerate(algo_data):
            metric_data[algo_name].append(getattr(
                one_data,
                metric_attr_name))  # get the attribute by the attr name

    return metric_data


def draw_cdf(cdf_data: dict[str, list[float]], metric_name: str):
    markers = ["*", "v", "+", "x", "d", "1"]
    marker_idx = 0

    plt.figure()
    for algo_name, algo_data in cdf_data.items():
        sorted_data = np.sort(algo_data)
        cumulative_prob = np.arange(1, len(algo_data) + 1) / len(algo_data)
        plt.plot(sorted_data,
                 cumulative_prob,
                 marker=markers[marker_idx],
                 markevery=(0, 20),
                 label=algo_name)
        marker_idx += 1
        if marker_idx >= len(markers):
            marker_idx = 0
    plt.title('Cumulative Distribution Function (CDF)')
    plt.xlabel(metric_name)
    plt.ylabel('Cumulative Probability')
    plt.grid(True)
    plt.legend()
    plt.show()


def main():
    all_data: dict[str, list[data_types.ResultData]] = dict()
    for _, algo_name in enumerate(ALGO_NAMES):  # initialize the dict
        all_data[algo_name] = []

    # load all data
    for i in range(REPEAT_COUNT):
        for algo_name in ALGO_NAMES:
            # use glob to get the paths of all csv files in this folder
            csv_files_pattern = "data/repeat{}/{}/*.csv".format(
                i + 1, algo_name)
            csv_file_names = glob.glob(csv_files_pattern)
            # read every csv files in this folder
            for _, file_name in enumerate(csv_file_names):
                all_data[algo_name].extend(csv_operation.read_csv(file_name))

    # complement the data for rejected applications
    comp_count_per_algo = REPEAT_COUNT * DEVICE_COUNT * APP_COUNT * REQ_COUNT_PER_APP
    data_for_rej_apps = calc_values_for_rejected(all_data)

    print(f"Every algorithm should have {comp_count_per_algo} items of data.")
    for _, algo_name in enumerate(ALGO_NAMES):
        algo_data = all_data[algo_name]
        num_to_comp = comp_count_per_algo - len(algo_data)

        print(
            f"Algorithm: {algo_name} has {len(algo_data)} items, so complement {num_to_comp} items to it."
        )

        all_data[algo_name].extend(
            [data_for_rej_apps for i in range(num_to_comp)])

    # # print the complemented results
    # for algo_name, algo_data in all_data.items():
    #     print(f"Algorithm: {algo_name} has {len(algo_data)} items.")
    #     for _, item in enumerate(algo_data):
    #         print(item)
    #     print()
    #     print()

    # draw CDF charts for every metric
    for attr_name, metric_name in ATTR_TO_METRIC.items():
        metric_data = make_cdf_data(all_data, attr_name)
        draw_cdf(metric_data, metric_name)


if __name__ == "__main__":
    main()
