import matplotlib.pyplot as plt
import numpy as np
import csv

ALGO_NAMES = [
    "BERand", "Amaga", "Ampga", "Diktyoga", "Mcssga"
]  # the names of all algorithms to draw. We do not draw "CompRand", because it may get unusable scheduling schemes.

APP_COUNTS = [70, 85, 100, 115, 130]

DATA_FILE_NAME_FMT = "usable_acceptance_rate_{app_count}.csv"


def main():
    # create the variable to store data
    # inner: app counts, outer: algorithms
    all_data: list[list[float]] = [[0] * (len(APP_COUNTS))
                                   for _ in range(0, len(ALGO_NAMES))]

    for app_idx, app_count in enumerate(APP_COUNTS):
        with open(DATA_FILE_NAME_FMT.format(app_count=app_count),
                  'r') as csv_file:
            csv_reader = csv.reader(csv_file, delimiter=",")

            next(csv_reader)  # skip the header row
            next(csv_reader)  # skip the row of algorithm "CompRand"
            algo_idx = 0
            for row in csv_reader:  # read total accepance rate of all priorities each row
                all_data[algo_idx][app_idx] = float(row[9])
                algo_idx += 1

    print("data read from file:")
    for i, one_pri_data in enumerate(all_data):
        print(i, one_pri_data)
    print()

    # width of each bar
    bar_inner_width = 0.1
    bar_outer_width = 0.15

    x_pos = np.arange(
        len(APP_COUNTS))  # generate the position of each group of bars

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
    hatches = [
        '///', '\\\\\\', '...', '---', 'xxx', '|||', '+++', 'ooo', '))O', '***'
    ]

    plt.rcParams["font.family"] = "Times New Roman"
    plt.rcParams['font.size'] = 15

    for i, one_algo_data in enumerate(all_data):
        plt.bar(x_pos + first_bar_offset + i * bar_outer_width,
                all_data[i],
                bar_inner_width,
                edgecolor='black',
                label=ALGO_NAMES[i],
                hatch=hatches[i])

    plt.xlabel("Number of applications")
    plt.ylabel("Acceptance rate")
    plt.title("Acceptance rate of applications with all priorities")
    plt.xticks(x_pos, APP_COUNTS)
    plt.legend()
    plt.grid(True, axis='y')
    plt.show()


if __name__ == "__main__":
    main()
