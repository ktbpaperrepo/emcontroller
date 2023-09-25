import matplotlib.pyplot as plt
import numpy as np
import csv

ALGO_NAMES = [
    "CompRand", "BERand", "Amaga", "Ampga", "Mcssga"
]  # the names of all algorithms to be evaluated in this experiment

ALGO_NAMES_TO_DRAW = [
    "BERand", "Amaga", "Ampga", "Mcssga"
]  # the names of all algorithms to draw. We do not draw "CompRand", because it may get unusable scheduling schemes.

#  the minimum and maximum possible priorities
MIN_PRI = 1
MAX_PRI = 10

DATA_FILE_NAME = "usable_acceptance_rate.csv"


def main():
    # create the variable to store data
    all_pri_data: list[dict[str, float]] = [
        dict() for _ in range(MIN_PRI, MAX_PRI + 1)
    ]

    # read data from the file
    with open(DATA_FILE_NAME, 'r') as csv_file:
        csv_reader = csv.reader(csv_file, delimiter=",")

        next(csv_reader)
        for row in csv_reader:  # read each row
            algo_name = row[0]
            pri = MIN_PRI
            col_idx = 11  # this column is the start of priority-separated data

            while pri <= MAX_PRI:
                all_pri_data[pri - 1][algo_name] = float(row[col_idx])
                pri += 1
                col_idx += 1

    print("data read from file:")
    for i, one_pri_data in enumerate(all_pri_data):
        print(i + 1, one_pri_data)
    print()

    # convert data for drawing
    data_for_drawing: list[list[float]] = [
        [0] * (MAX_PRI - MIN_PRI + 1)
        for _ in range(0, len(ALGO_NAMES_TO_DRAW))
    ]
    for i, one_pri_data in enumerate(all_pri_data):
        for j, algo_name in enumerate(ALGO_NAMES_TO_DRAW):
            data_for_drawing[j][i] = one_pri_data[algo_name]

    # draw the bar chart to compare the app acceptance rate of every priority
    print("data for drawing:")
    for i, one_algo_data in enumerate(data_for_drawing):
        print(ALGO_NAMES_TO_DRAW[i], one_algo_data)
    print()

    # width of each bar
    bar_inner_width = 0.12
    bar_outer_width = 0.17

    x_pos = np.arange(MAX_PRI - MIN_PRI +
                      1)  # generate the position of each group of bars

    # calculate the offsets of all bars
    bar_count = len(data_for_drawing)
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

    for i, one_algo_data in enumerate(data_for_drawing):
        plt.bar(x_pos + first_bar_offset + i * bar_outer_width,
                data_for_drawing[i],
                bar_inner_width,
                edgecolor='black',
                label=ALGO_NAMES_TO_DRAW[i],
                hatch=hatches[i])

    plt.xlabel("Application priority")
    plt.ylabel("Acceptance rate")
    plt.title("Acceptance rate of applications with different priorities")
    plt.xticks(x_pos, range(MIN_PRI, MAX_PRI + 1))
    plt.legend()
    plt.grid(True, axis='y')
    plt.show()


if __name__ == "__main__":
    main()
