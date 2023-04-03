'use strict';

// For safety, we do not allow users to delete VMs too frequently
let deleteVMLock = false;

function lockDeleteVM() {
    deleteVMLock = true;
}

function unlockDeleteVM() {
    deleteVMLock = false;
    location.reload() // after deleting, refresh the page
}

// original html does not support to send PUT or DELETE request
function deleteVM(cloudName, vmID, statusID) {
    if (deleteVMLock) {
        console.log("Another deleting is executing, please try again after a few seconds");
        return;
    }
    lockDeleteVM();
    let vmStatus = document.getElementById(statusID);
    vmStatus.innerText = "Deleting";
    let xmlhttp = new XMLHttpRequest();
    xmlhttp.open("DELETE", `/cloud/${cloudName}/vm/${vmID}`);
    xmlhttp.send();
    console.log("delete vm %s in cloud %s request has been sent", vmID, cloudName);
    xmlhttp.onreadystatechange = function(){
        if(this.readyState==4) {
            if (this.status==200) {
                console.log("delete vm %s in cloud %s response: %s", vmID, cloudName, xmlhttp.responseText);
                // this Delete VM API will block until the VM is completely deleted, so 2 seconds of waiting is enough
                setTimeout(unlockDeleteVM, 2000);
            } else {
                vmStatus.innerText = `Delete Error:\r\nreadyState: ${this.readyState}\r\nstatus: ${this.status}`;
            }
        }
        console.log("onreadystatechange this.readyState: %O, this.status: %O", this.readyState, this.status);
    }
}

function whileCreatingVM() {
    let newVmName = document.getElementById("newVmName");
    let newVmVCpu = document.getElementById("newVmVCpu");
    let newVmRam = document.getElementById("newVmRam");
    let newVmStorage = document.getElementById("newVmStorage");
    let createButton = document.getElementById("createNewVm");

    newVmName.setAttribute("readonly", "readonly");
    newVmVCpu.setAttribute("readonly", "readonly");
    newVmRam.setAttribute("readonly", "readonly");
    newVmStorage.setAttribute("readonly", "readonly");
    createButton.setAttribute("disabled", "disabled");
    createButton.insertAdjacentHTML('afterend',`<p>Creating VM ${newVmName.value}, please wait ...</p>`);
}