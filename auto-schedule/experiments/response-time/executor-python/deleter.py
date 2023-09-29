import time

import http_api
import data_types

EXPT_APP_NAME_PREFIX = "expt-app-"
AUTO_SCHED_VM_NAME_PREFIX = "auto-sched-"


# wait until a desired state, checking with the input function checkFunc
def poll(maxCheckCount: int, checkIntervalSecond: int, checkFunc):
    checkCount = 1
    while checkCount < maxCheckCount:
        print("The {} time check.".format(checkCount))
        checkCount += 1
        success, err_str = checkFunc()
        if success:
            return True, ""
        if len(err_str) > 0:
            return False, err_str
        time.sleep(checkIntervalSecond)

    return False, "The desired state did not appear after {} times of check with interval {} seconds.".format(
        checkCount, checkIntervalSecond)


# check whether Auto-scheduling Virtual Machines are already cleaned by the garbage collection of multi-cloud manager
def check_as_vms_deleted():
    all_vms = http_api.get_all_vms()
    for idx, vm in enumerate(all_vms):
        if vm.name.startswith(AUTO_SCHED_VM_NAME_PREFIX):
            print("Virtual Machine {} has not been deleted.".format(vm.name))
            return False, ""
    return True, ""


# check whether Auto-scheduling Kubernetes nodes are already cleaned by the garbage collection of multi-cloud manager
def check_as_nodes_deleted():
    all_k8s_nodes = http_api.get_k8s_nodes()
    for idx, node in enumerate(all_k8s_nodes):
        if node.name.startswith(AUTO_SCHED_VM_NAME_PREFIX):
            print("Kubernetes nodes {} has not been deleted.".format(
                node.name))
            return False, ""
    return True, ""


def main():

    # this is the data structure of the applications
    apps: list[data_types.AppInfo] = []

    apps = http_api.get_all_apps()  # get all applications

    # put the name of applications in an array, which can be the json body of the HTTP request to delete applications via multi-cloud manager
    app_names: list[str] = []
    for idx, app in enumerate(apps):
        if app.appName.startswith(EXPT_APP_NAME_PREFIX):
            app_names.append(app.appName)

    print("send API to multi-cloud manager to delete applications:", app_names)
    http_api.del_apps(app_names)

    # After delete applications, we should wait until the GC of multi-cloud manager clean up the auto-scheduling Kubernetes nodes and Virtual Machines.
    success, err_str = poll(100, 60, check_as_vms_deleted)
    if success:
        print(
            "Auto-scheduling Virtual Machines are already cleaned by the garbage collection of multi-cloud manager."
        )
    else:
        raise Exception(
            "Failed to wait for Virtual Machines cleaned by the garbage collection of multi-cloud manager, the error message is: {}"
            .format(err_str))
    # Multi-cloud manager GC cleans K8s nodes firstly and then VMs, so in principle we only need to check VMs being deleted, but for safety, we also do a check for K8s nodes after VMs are cleaned.
    success, err_str = poll(100, 60, check_as_nodes_deleted)
    if success:
        print(
            "Auto-scheduling Kubernetes nodes are already cleaned by the garbage collection of multi-cloud manager."
        )
    else:
        raise Exception(
            "Failed to wait for Kubernetes nodes cleaned by the garbage collection of multi-cloud manager, the error message is: {}"
            .format(err_str))


if __name__ == "__main__":
    main()
