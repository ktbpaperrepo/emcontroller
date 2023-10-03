import matplotlib.pyplot as plt
import numpy as np
import csv

import common

ALGO_NAMES = [
    "CompRand", "BERand", "Amaga", "Ampga", "Diktyoga", "Mcssga"
]  # the names of all algorithms to be evaluated in this experiment


def main():
    # create the variable to store data
    # inner: app counts, outer: algorithms
    all_data: list[list[float]] = [[0] * (len(common.APP_COUNTS))
                                   for _ in range(0, len(ALGO_NAMES))]

    for app_count_idx, app_count in enumerate(common.APP_COUNTS):
        with open(common.DATA_FILE_NAME_FMT.format(app_count=app_count),
                  'r') as csv_file:
            csv_reader = csv.reader(csv_file, delimiter=",")

            next(csv_reader)
            algo_idx = 0
            for row in csv_reader:  # read maximum scheduling time in each row
                all_data[algo_idx][app_count_idx] = float(row[1])
                algo_idx += 1

    print("data read from file:")
    for i, one_pri_data in enumerate(all_data):
        print(i, one_pri_data)
    print()

    # width of each bar
    bar_inner_width = 0.1
    bar_outer_width = 0.15

    x_pos = np.arange(len(
        common.APP_COUNTS))  # generate the position of each group of bars

    # calculate the offsets of all bars
    bar_count = len(all_data)
    first_bar_offset: float = 0
    if bar_count % 2 == 0:
        first_bar_offset = -bar_outer_width / 2 - bar_outer_width * (
            bar_count / 2 - 1)
    else:
        first_bar_offset = -bar_outer_width - bar_outer_width * (
            (bar_count - 1) / 2 - 1)

    # draw bars
    hatches = common.HATCHES

    plt.figure(figsize=common.FIG_SIZE)
    plt.rcParams["font.family"] = common.FONT_FAMILY
    plt.rcParams['font.size'] = common.FONT_SIZE

    for i, one_algo_data in enumerate(all_data):
        plt.bar(x_pos + first_bar_offset + i * bar_outer_width,
                all_data[i],
                bar_inner_width,
                edgecolor='black',
                label=ALGO_NAMES[i],
                hatch=hatches[i])

    plt.xlabel("Number of applications")
    plt.ylabel("Maximum scheduling time (s)")
    plt.title("Time of scheduling different numbers of applications")
    plt.xticks(x_pos, common.APP_COUNTS)
    plt.legend()
    plt.grid(True, axis='y')
    plt.show()


if __name__ == "__main__":
    main()
