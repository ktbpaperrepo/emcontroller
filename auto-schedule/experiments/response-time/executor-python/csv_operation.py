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
