import csv

import data_types


def write_csv(csv_file_name: str, results: list[data_types.ResultData]):
    with open(csv_file_name, 'w') as csv_file:
        writer = csv.writer(csv_file, delimiter=",")
        writer.writerow([
            "app_name", "priority", "resp_time", "resp_time_in_clouds",
            "pri_wei_resp_time", "pri_wei_resp_time_in_clouds"
        ])

        for i, result in enumerate(results):
            writer.writerow([
                result.app_name, result.priority, result.resp_time,
                result.resp_time_in_clouds, result.pri_wei_resp_time,
                result.pri_wei_resp_time_in_clouds
            ])


def read_csv(csv_file_name: str) -> list[data_types.ResultData]:
    results: list[data_types.ResultData] = []

    with open(csv_file_name, 'r') as csv_file:
        csv_reader = csv.reader(csv_file, delimiter=",")

        next(csv_reader)  # skip the first row of the csv file.
        for row in csv_reader:
            results.append(
                data_types.ResultData(app_name=row[0],
                                      priority=int(row[1]),
                                      resp_time=float(row[2]),
                                      resp_time_in_clouds=float(row[3]),
                                      pri_wei_resp_time=float(row[4]),
                                      pri_wei_resp_time_in_clouds=float(
                                          row[5])))

    return results
