from concurrent import futures
import multiprocessing

import http_api
import data_types
import csv_operation
import other_tools

EXPT_APP_NAME_PREFIX = "expt-app-"
num_cores = multiprocessing.cpu_count()


def main():

    # this is the data structure of the applications
    apps: list[data_types.AppInfo] = []

    apps = http_api.get_all_apps()  # get all applications

    results: list[data_types.ResultData] = []  # we will put results here

    # call the applications in parallel
    process_pool = futures.ProcessPoolExecutor(max_workers=num_cores)
    future_list = []

    for idx, app in enumerate(apps):
        # if idx > 5: # for debug
        #     break

        if not app.appName.startswith(EXPT_APP_NAME_PREFIX):
            continue
        # print("app {}: {}".format(idx, app))
        print(
            "Submit request to call App No. {}, Name: {}, Priority: {},NodePort endpoint: {}:{}."
            .format(idx, app.appName, app.priority, app.nodePortIP[0],
                    app.nodePort[0]))

        # submit
        future_list.append(process_pool.submit(http_api.call_app, app))

    print('Start to wait for responses')
    process_pool.shutdown(wait=True)
    print('All responses are received!')

    # put the response data into results
    for future in futures.as_completed(future_list):
        results.append(future.result())

    # print(results)
    data_file_name = other_tools.gen_data_file_name()
    csv_operation.write_csv(data_file_name, results)


if __name__ == "__main__":
    main()
