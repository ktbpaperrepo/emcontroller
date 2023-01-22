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
    let appStatus = document.getElementById(statusID);
    appStatus.innerText = "Deleting";
    let xmlhttp = new XMLHttpRequest();
    xmlhttp.open("DELETE", `/cloud/${cloudName}/vm/${vmID}`);
    xmlhttp.send();
    console.log("delete vm %s request has been sent", vmID);
    xmlhttp.onreadystatechange = function(){
        if(this.readyState==4 && this.status==200) {
            console.log("delete vm %s response: %s", vmID, xmlhttp.responseText);
            // reserve 10s for deleting
            setTimeout(unlockDeleteVM, 10000);
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
