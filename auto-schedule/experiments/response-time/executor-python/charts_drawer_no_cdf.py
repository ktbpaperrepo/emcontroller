import glob
import json
import numpy as np
import matplotlib.pyplot as plt

import data_types
import csv_operation
import charts_drawer

AVG_NON_PRI_ATTR_TO_METRIC: dict[str, str] = {
    "resp_time": "Average response time (ms)",
    "resp_time_in_clouds": "Average response time in clouds (ms)",
}


# sort the input 2 mapped lists by priority
def mapped_sort_pri(data_list1: list[data_types.ResultData],
                    data_list2: list[data_types.ResultData]):
    sorted_list1: list[data_types.ResultData] = []
    sorted_list2: list[data_types.ResultData] = []

    # traverse all priorities from small to big
    for pri in range(charts_drawer.MIN_PRI, charts_drawer.MAX_PRI + 1):
        for idx, _ in enumerate(data_list1):
            if data_list1[idx].priority == pri:
                sorted_list1.append(data_list1[idx])
                sorted_list2.append(data_list2[idx])

    if len(sorted_list1) != len(data_list1) or len(sorted_list2) != len(
            data_list2):
        raise Exception(
            "len(sorted_list1): {}, len(data_list1): {}, len(sorted_list2): {}, len(data_list2): {}"
            .format(len(sorted_list1), len(data_list1), len(sorted_list2),
                    len(data_list2)))

    validate_one_sort_list(sorted_list1)
    validate_one_sort_list(sorted_list2)

    return sorted_list1, sorted_list2


# validate that the list items of the 2 algorithms can be mapped.
def validate_both_list(
    both_acc_list: dict[str, dict[str, dict[str,
                                            list[data_types.ResultData]]]]):
    print("start validating")
    for j, _ in enumerate(charts_drawer.ALGO_NAMES):
        algo1 = charts_drawer.ALGO_NAMES[j]
        for k in range(j + 1, len(charts_drawer.ALGO_NAMES)):
            algo2 = charts_drawer.ALGO_NAMES[k]

            if len(both_acc_list[algo1][algo2][algo1]) != len(
                    both_acc_list[algo1][algo2][algo1]):
                raise Exception("length not equal")

            for idx, _ in enumerate(both_acc_list[algo1][algo2][algo1]):
                if both_acc_list[algo1][algo2][algo1][
                        idx].app_name != both_acc_list[algo1][algo2][algo2][
                            idx].app_name or both_acc_list[algo1][algo2][
                                algo1][idx].priority != both_acc_list[algo1][
                                    algo2][algo2][idx].priority:
                    raise Exception(
                        "{}, {}, {}, not equal, {}, {}, {}.".format(
                            algo1,
                            both_acc_list[algo1][algo2][algo1][idx].app_name,
                            both_acc_list[algo1][algo2][algo1][idx].priority,
                            algo2,
                            both_acc_list[algo1][algo2][algo2][idx].app_name,
                            both_acc_list[algo1][algo2][algo2][idx].priority))


def validate_one_sort_list(data_list: list[data_types.ResultData]):
    current_pri = charts_drawer.MIN_PRI
    for idx, _ in enumerate(data_list):
        pri = data_list[idx].priority
        if pri < current_pri:
            raise Exception(
                "priority {} is smaller than the current {}.".format(
                    pri, current_pri))
        if pri > current_pri:
            current_pri = pri
            # print("current is changed to {}".format(current_pri))


# validate that the sorted list items of the 2 algorithms can be mapped.
def validate_sorted_both_list(
    sorted_both_acc_list: dict[str, dict[str,
                                         dict[str,
                                              list[data_types.ResultData]]]]):

    validate_both_list(sorted_both_acc_list)

    print("start sorted validating")

    for j, _ in enumerate(charts_drawer.ALGO_NAMES):
        algo1 = charts_drawer.ALGO_NAMES[j]
        for k in range(j + 1, len(charts_drawer.ALGO_NAMES)):
            algo2 = charts_drawer.ALGO_NAMES[k]
            try:
                validate_one_sort_list(
                    sorted_both_acc_list[algo1][algo2][algo1])
            except Exception as e:
                raise Exception("{}, {}, {}, error: {}.".format(
                    algo1, algo2, algo1, e))


# filter the applications accepted by all algorithms
def filter_all_accepted(
    all_data: dict[str, dict[str, list[data_types.ResultData]]],
    all_app_name_to_pri: dict[str, int], algo_names_to_compare: list[str]
) -> dict[str, dict[str, list[data_types.ResultData]]]:

    all_accepted_data: dict[str, dict[str,
                                      list[data_types.ResultData]]] = dict()
    for _, algo_name in enumerate(
            algo_names_to_compare):  # initialize the dict
        all_accepted_data[algo_name] = dict()

    for app_name, _ in all_app_name_to_pri.items():
        rejected = False
        for _, algo_name in enumerate(algo_names_to_compare):
            if len(all_data[algo_name][app_name]) == 0:
                rejected = True
        if not rejected:
            for _, algo_name in enumerate(algo_names_to_compare):
                all_accepted_data[algo_name][app_name] = all_data[algo_name][
                    app_name]

    return all_accepted_data


# extract the data of the applications accepted by both of every 2 algorithms
def extract_both_accepted_data(
    all_data: dict[str, dict[str, list[data_types.ResultData]]],
    all_app_name_to_pri: dict[str, int]
) -> dict[str, dict[str, dict[str, dict[str, list[data_types.ResultData]]]]]:

    # init the return
    algos_to_both_acc_data: dict[str, dict[str, dict[str, dict[
        str, list[data_types.ResultData]]]]] = dict()
    for i, _ in enumerate(charts_drawer.ALGO_NAMES):
        algos_to_both_acc_data[charts_drawer.ALGO_NAMES[i]] = dict()

    # compare every 2 algorithms
    for i, _ in enumerate(charts_drawer.ALGO_NAMES):
        for j in range(i + 1, len(charts_drawer.ALGO_NAMES)):
            algos_to_cmp = [
                charts_drawer.ALGO_NAMES[i], charts_drawer.ALGO_NAMES[j]
            ]
            # only draw charts for the applications accepted by both of the selected algorithms
            app_data_all_accepted = filter_all_accepted(
                all_data, all_app_name_to_pri, algos_to_cmp)
            # # draw
            # draw_cdf_every_metric(
            #     app_data_all_accepted, 3,
            #     "Repeat {}. Apps accepted by both {} and {}".format(
            #         repeat_idx, ALGO_NAMES[i], ALGO_NAMES[j]))
            # save data
            algos_to_both_acc_data[charts_drawer.ALGO_NAMES[i]][
                charts_drawer.ALGO_NAMES[j]] = app_data_all_accepted

    return algos_to_both_acc_data


# draw charts for metrics not weighted by priorities, because we will show the effects of priorities in another way
def draw_charts_no_pri_metrics(
        algo_to_data: dict[str, list[data_types.ResultData]]):
    for attr_name, metric_name in AVG_NON_PRI_ATTR_TO_METRIC.items():
        draw_dots_chart(algo_to_data, attr_name, metric_name)


# draw dots in a charts for the input data
def draw_dots_chart(algo_to_data: dict[str, list[data_types.ResultData]],
                    attr_name: str, metric_name: str):

    markers = ["*", "v", "+", "x", "d", "1"]
    colors = [
        '#1f77b4', '#ff7f0e', '#2ca02c', '#d62728', '#9467bd', '#8c564b',
        '#e377c2', '#7f7f7f', '#bcbd22', '#17becf'
    ]
    marker_color_idx = 0

    dot_size = 22

    plt.figure(figsize=(10, 6))

    # # trigger core fonts for PDF backend
    # plt.rcParams["pdf.use14corefonts"] = True
    # # trigger core fonts for PS backend
    # plt.rcParams["ps.useafm"] = True

    plt.rcParams["font.family"] = "Times New Roman"
    plt.rcParams['font.size'] = 17

    # draw dots
    algo_names = []
    x_values = []  # for the following plotting
    algo_y_values = []  # for the following plotting
    for algo_name, algo_data in algo_to_data.items():
        this_x_values, this_y_values = draw_dots(algo_name, algo_data,
                                                 attr_name,
                                                 markers[marker_color_idx],
                                                 colors[marker_color_idx],
                                                 dot_size)
        marker_color_idx += 1
        algo_names.append(algo_name)
        x_values = this_x_values
        algo_y_values.append(this_y_values)

    # the 2 horizontal lines to show which algorithm has a shorter response time

    # draw the 0 line to avoid confusion
    plt.axhline(y=0, color='black', linestyle='-', linewidth=1.3, zorder=1)
    y_bottom_lim = -plt.ylim()[1] * 0.3
    algo_hline_pos = [-plt.ylim()[1] * 0.1, -(plt.ylim()[1] * 0.2)]
    plt.axhline(y=algo_hline_pos[0],
                color='black',
                linestyle='--',
                linewidth=1,
                zorder=1)

    plt.axhline(y=algo_hline_pos[1],
                color='black',
                linestyle='--',
                linewidth=1,
                zorder=1)

    plt.ylim(bottom=y_bottom_lim)

    # draw dots on the lines to show the shorter response time
    algo_better_nums = [0, 0]
    for i, _ in enumerate(x_values):
        x = x_values[i]
        y1 = algo_y_values[0][i]
        y2 = algo_y_values[1][i]
        algo1 = algo_names[0]
        algo2 = algo_names[1]
        marker1 = markers[0]
        marker2 = markers[1]
        color1 = colors[0]
        color2 = colors[1]

        if y1 < y2:
            algo_better_nums[0] += 1
            plt.scatter(x,
                        algo_hline_pos[0],
                        marker=marker1,
                        color=color1,
                        s=dot_size)
        elif y1 > y2:
            algo_better_nums[1] += 1
            plt.scatter(x,
                        algo_hline_pos[1],
                        marker=marker2,
                        color=color2,
                        s=dot_size)

    # write text to explain the horizontal lines
    plt.text(plt.xlim()[0] - (plt.xlim()[1] - plt.xlim()[0]) * 0.003,
             algo_hline_pos[0],
             "{} faster".format(algo_names[0]),
             ha='right',
             va='top',
             size=15,
             rotation=45)
    plt.text(plt.xlim()[1] + (plt.xlim()[1] - plt.xlim()[0]) * 0.005,
             algo_hline_pos[0],
             "{} apps".format(algo_better_nums[0]),
             ha='left',
             va='center',
             size=15)
    plt.text(plt.xlim()[0] - (plt.xlim()[1] - plt.xlim()[0]) * 0.003,
             algo_hline_pos[1],
             "{} faster".format(algo_names[1]),
             ha='right',
             va='top',
             size=15,
             rotation=45)
    plt.text(plt.xlim()[1] + (plt.xlim()[1] - plt.xlim()[0]) * 0.005,
             algo_hline_pos[1],
             "{} apps".format(algo_better_nums[1]),
             ha='left',
             va='center',
             size=15)

    # adjust the ticks after finishing the size, before using plt.ylim(), because ticks will change plt.ylim()

    plt.xticks([])
    # remove the negative values from the y_ticks
    existing_y_ticks, _ = plt.yticks()
    first_positive_idx = 0
    for idx, tick in enumerate(existing_y_ticks):
        if tick >= 0:
            first_positive_idx = idx
            break
    plt.yticks(existing_y_ticks[first_positive_idx:])

    plt.ylabel(metric_name, loc="top")
    plt.xlabel('Applications accepted by both {} and {}'.format(
        algo_names[0], algo_names[1]))

    # the vertical lines to split different priorities
    pri_change_points = []
    priorities = []
    for idx, _ in enumerate(algo_to_data[algo_names[0]]):
        if idx == 0:
            priorities.append(algo_to_data[algo_names[0]][idx].priority)
            continue
        this_data = algo_to_data[algo_names[0]][idx]
        last_data = algo_to_data[algo_names[0]][idx - 1]
        if this_data.priority != last_data.priority:
            pri_change_points.append(idx - 0.5)
            priorities.append(algo_to_data[algo_names[0]][idx].priority)

    pri_poses = []
    for i, point in enumerate(pri_change_points):
        if i == 0:
            pri_poses.append(plt.xlim()[0] + (point - plt.xlim()[0]) / 2)
        else:
            pri_poses.append(pri_change_points[i - 1] +
                             (point - pri_change_points[i - 1]) / 2)
    if len(pri_change_points) == 0:
        pri_poses.append(plt.xlim()[0] + (plt.xlim()[1] - plt.xlim()[0]) / 2)
    else:
        pri_poses.append(pri_change_points[len(pri_change_points) - 1] +
                         (plt.xlim()[1] -
                          pri_change_points[len(pri_change_points) - 1]) / 2)

    for _, pri_change_point in enumerate(pri_change_points):
        plt.axvline(x=pri_change_point,
                    color='black',
                    linestyle='-',
                    linewidth=0.5,
                    zorder=1)

    # the text to explain the priorities
    for i, _ in enumerate(pri_poses):
        plt.text(pri_poses[i],
                 plt.ylim()[1] + (plt.ylim()[1] - plt.ylim()[0]) * 0.01,
                 priorities[i],
                 ha='center')

    plt.text(plt.xlim()[0] + (plt.xlim()[1] - plt.xlim()[0]) / 2,
             plt.ylim()[1] + (plt.ylim()[1] - plt.ylim()[0]) * 0.07,
             "Application priorities",
             ha='center')

    plt.legend()
    plt.show()


# draw dots for one group of data
def draw_dots(algo_name: str, algo_data: list[data_types.ResultData],
              attr_name: str, marker: str, color: str, size: int):
    # generate the values to plot
    y_values: list[float] = []
    for _, one_data in enumerate(algo_data):
        y_values.append(getattr(
            one_data, attr_name))  # get the attribute by the attr name

    # we do not need x values, so x_values are the indexes of y_values
    x_values = range(len(y_values))
    plt.scatter(x_values,
                y_values,
                marker=marker,
                color=color,
                s=size,
                zorder=5,
                label=algo_name)

    return x_values, y_values,


def main():

    # initializa this variable to save data to merge the applications accepted by both of every 2 algorithms
    data_acc_both_algos: dict[str, dict[str, dict[
        str, list[data_types.ResultData]]]] = dict()
    for i, _ in enumerate(charts_drawer.ALGO_NAMES):
        data_acc_both_algos[charts_drawer.ALGO_NAMES[i]] = dict()
        for j in range(i + 1, len(charts_drawer.ALGO_NAMES)):
            data_acc_both_algos[charts_drawer.ALGO_NAMES[i]][
                charts_drawer.ALGO_NAMES[j]] = dict()
            data_acc_both_algos[charts_drawer.ALGO_NAMES[i]][
                charts_drawer.ALGO_NAMES[j]][charts_drawer.ALGO_NAMES[i]] = []
            data_acc_both_algos[charts_drawer.ALGO_NAMES[i]][
                charts_drawer.ALGO_NAMES[j]][charts_drawer.ALGO_NAMES[j]] = []

    # initialize this variable to store the data for ploting
    avg_all_repeats_both_list: dict[str, dict[str, dict[
        str, list[data_types.ResultData]]]] = dict()
    for j, _ in enumerate(charts_drawer.ALGO_NAMES):
        algo1 = charts_drawer.ALGO_NAMES[j]
        avg_all_repeats_both_list[algo1] = dict()
        for k in range(j + 1, len(charts_drawer.ALGO_NAMES)):
            algo2 = charts_drawer.ALGO_NAMES[k]
            avg_all_repeats_both_list[algo1][algo2] = dict()
            avg_all_repeats_both_list[algo1][algo2][algo1] = []
            avg_all_repeats_both_list[algo1][algo2][algo2] = []

    # load all data
    for i in range(charts_drawer.REPEAT_COUNT):
        # in a repeat, the apps are the same for all algorithms
        app_name_to_pri: dict[str, int] = dict()

        # from the app deploy request json to read all applications including rejected ones
        with open("data/repeat{}/request_applications.json".format(
                i + 1)) as json_file:
            app_reqs = json.load(json_file)
            for _, app_req in enumerate(app_reqs):
                app_name_to_pri[app_req["name"]] = app_req["priority"]

        # to save data of this repeat
        all_data_this_repeat: dict[str,
                                   dict[str,
                                        list[data_types.ResultData]]] = dict()
        for _, algo_name in enumerate(
                charts_drawer.ALGO_NAMES):  # initialize the dict
            all_data_this_repeat[algo_name] = dict()
            for app_name, _ in app_name_to_pri.items():
                all_data_this_repeat[algo_name][app_name] = []

        # read the data from csv files
        for algo_name in charts_drawer.ALGO_NAMES:
            # use glob to get the paths of all csv files in this folder
            csv_files_pattern = "data/repeat{}/{}/*.csv".format(
                i + 1, algo_name)
            csv_file_names = glob.glob(csv_files_pattern)
            # read every csv files in this folder
            for _, file_name in enumerate(csv_file_names):
                data_in_file = csv_operation.read_csv(file_name)
                # to save data of this repeat
                for _, one_item in enumerate(data_in_file):
                    all_data_this_repeat[algo_name][one_item.app_name].append(
                        one_item)

        # # print for debug
        # for algo_name, algo_data in all_data_this_repeat.items():
        #     for app_name, app_data in algo_data.items():
        #         print("Repeat {}, algorithm {}, app {}, data: {}".format(
        #             i + 1, algo_name, app_name, app_data))

        this_repeat_algos_to_both_acc_data = extract_both_accepted_data(
            all_data_this_repeat, app_name_to_pri)

        # calculate the average value
        avg_this_repeat_both: dict[str, dict[str, dict[str, dict[
            str, data_types.ResultData]]]] = dict()
        for algo1, algo1_dict in this_repeat_algos_to_both_acc_data.items():
            avg_this_repeat_both[algo1] = dict()
            for algo2, both_data in algo1_dict.items():
                avg_this_repeat_both[algo1][algo2] = dict()
                avg_this_repeat_both[algo1][algo2][algo1] = dict()
                avg_this_repeat_both[algo1][algo2][algo2] = dict()

                for app_name, app_data in both_data[algo1].items():
                    avg_this_repeat_both[algo1][algo2][algo1][
                        app_name] = data_types.calc_rd_avg(app_data)
                for app_name, app_data in both_data[algo2].items():
                    avg_this_repeat_both[algo1][algo2][algo2][
                        app_name] = data_types.calc_rd_avg(app_data)

        # # print for debug
        # # merge data of all repeats
        # for j, _ in enumerate(charts_drawer.ALGO_NAMES):
        #     for k in range(j + 1, len(charts_drawer.ALGO_NAMES)):
        #         for app_name, app_data in this_repeat_algos_to_both_acc_data[
        #                 charts_drawer.ALGO_NAMES[j]][charts_drawer.ALGO_NAMES[
        #                     k]][charts_drawer.ALGO_NAMES[j]].items():
        #             if charts_drawer.ALGO_NAMES[
        #                     j] == "Diktyoga" and charts_drawer.ALGO_NAMES[
        #                         k] == "Mcssga" and app_name == "expt-app-99":
        #                 print(i + 1, charts_drawer.ALGO_NAMES[j],
        #                       charts_drawer.ALGO_NAMES[k], app_data)

        #         for app_name, app_data in avg_this_repeat_both[
        #                 charts_drawer.ALGO_NAMES[j]][charts_drawer.ALGO_NAMES[
        #                     k]][charts_drawer.ALGO_NAMES[j]].items():
        #             if charts_drawer.ALGO_NAMES[
        #                     j] == "Diktyoga" and charts_drawer.ALGO_NAMES[
        #                         k] == "Mcssga" and app_name == "expt-app-99":
        #                 print(i + 1, charts_drawer.ALGO_NAMES[j],
        #                       charts_drawer.ALGO_NAMES[k], app_data)

        #         for app_name, app_data in avg_this_repeat_both[
        #                 charts_drawer.ALGO_NAMES[j]][charts_drawer.ALGO_NAMES[
        #                     k]][charts_drawer.ALGO_NAMES[k]].items():
        #             if charts_drawer.ALGO_NAMES[
        #                     j] == "Diktyoga" and charts_drawer.ALGO_NAMES[
        #                         k] == "Mcssga" and app_name == "expt-app-99":
        #                 print(i + 1, charts_drawer.ALGO_NAMES[k],
        #                       charts_drawer.ALGO_NAMES[j], app_data)

        # put the data from dicts to lists
        avg_this_repeat_both_list: dict[str, dict[str, dict[
            str, list[data_types.ResultData]]]] = dict()

        for j, _ in enumerate(charts_drawer.ALGO_NAMES):
            algo1 = charts_drawer.ALGO_NAMES[j]

            avg_this_repeat_both_list[algo1] = dict()

            for k in range(j + 1, len(charts_drawer.ALGO_NAMES)):
                algo2 = charts_drawer.ALGO_NAMES[k]

                avg_this_repeat_both_list[algo1][algo2] = dict()

                avg_this_repeat_both_list[algo1][algo2][algo1] = []
                avg_this_repeat_both_list[algo1][algo2][algo2] = []
                for app_name, _ in avg_this_repeat_both[algo1][algo2][
                        algo1].items():
                    avg_this_repeat_both_list[algo1][algo2][algo1].append(
                        avg_this_repeat_both[algo1][algo2][algo1][app_name])
                    avg_this_repeat_both_list[algo1][algo2][algo2].append(
                        avg_this_repeat_both[algo1][algo2][algo2][app_name])

        # validate, not a function
        validate_both_list(avg_this_repeat_both_list)
        print(
            "Repeat {}, validate of \"avg_this_repeat_both_list\" is finished."
            .format(i + 1))

        # merge the data of this repeat to the total data
        for j, _ in enumerate(charts_drawer.ALGO_NAMES):
            algo1 = charts_drawer.ALGO_NAMES[j]
            for k in range(j + 1, len(charts_drawer.ALGO_NAMES)):
                algo2 = charts_drawer.ALGO_NAMES[k]
                avg_all_repeats_both_list[algo1][algo2][algo1].extend(
                    avg_this_repeat_both_list[algo1][algo2][algo1])
                avg_all_repeats_both_list[algo1][algo2][algo2].extend(
                    avg_this_repeat_both_list[algo1][algo2][algo2])

    validate_both_list(avg_all_repeats_both_list)
    print(
        "All repeats, validate of \"avg_all_repeats_both_list\" is finished.")

    # sort the applications by priority
    for j, _ in enumerate(charts_drawer.ALGO_NAMES):
        algo1 = charts_drawer.ALGO_NAMES[j]
        for k in range(j + 1, len(charts_drawer.ALGO_NAMES)):
            algo2 = charts_drawer.ALGO_NAMES[k]

            unsorted_list1 = avg_all_repeats_both_list[algo1][algo2][algo1]
            unsorted_list2 = avg_all_repeats_both_list[algo1][algo2][algo2]

            sorted_list1, sorted_list2 = mapped_sort_pri(
                unsorted_list1, unsorted_list2)

            avg_all_repeats_both_list[algo1][algo2][algo1] = sorted_list1
            avg_all_repeats_both_list[algo1][algo2][algo2] = sorted_list2

    validate_sorted_both_list(avg_all_repeats_both_list)
    print(
        "All repeats, validate of the sorted \"avg_all_repeats_both_list\" is finished."
    )

    # plotting
    for j, _ in enumerate(charts_drawer.ALGO_NAMES):
        algo1 = charts_drawer.ALGO_NAMES[j]
        for k in range(j + 1, len(charts_drawer.ALGO_NAMES)):
            algo2 = charts_drawer.ALGO_NAMES[k]
            if algo1 != charts_drawer.OUR_ALGO_NAME and algo2 != charts_drawer.OUR_ALGO_NAME:
                continue
            draw_charts_no_pri_metrics(avg_all_repeats_both_list[algo1][algo2])


if __name__ == "__main__":
    main()
