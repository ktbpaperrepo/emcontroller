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

// delete multiple VMs
function deleteBatchVMs() {
    if (deleteVMLock) {
        console.log("Another deleting is executing, please try again after a few seconds");
        return;
    }
    lockDeleteVM();

    let deleteSelectedButton = document.getElementById("deleteSelectedButton");
    deleteSelectedButton.insertAdjacentHTML('afterend',"<p id=\"textForDeleting\">Deleting selected VMs, please wait ...</p>");

    let vmsToDelete = [];
    let vmCheckboxes = document.getElementsByClassName("vmCheckbox");

    for (let i = 0; i < vmCheckboxes.length; i++) {
        if (vmCheckboxes[i].checked) {
            // get the row of the table
            let row = vmCheckboxes[i].parentNode.parentNode;

            // set the "Status" column as "Deleting"
            row.cells[11].innerText = "Deleting";

            // get the needed information to delete a VM
            let vmName = row.cells[3].textContent;
            let cloudName = row.cells[5].textContent;
            let vmId = row.cells[6].textContent;

            // make the json body for the request
            vmsToDelete.push({
                id: vmId,
                name: vmName,
                cloud: cloudName,
            });

        }
    }

    // send http request to delete VMs
    let resp = fetch("/vm",{
        method: "DELETE",
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify(vmsToDelete)
    })

    // The return type of fetch is Promise and the Promise can only be accessed in its then.
    // I need to access its status code and text.
    // - The status code can only be accessed in resp.then().
    // - The text can only be accessed in resp.text().then().
    // Therefore, I use a 2-layer then() to both access them.
    resp.then(response => {
        response.text().then(text => {
            if (response.status >= 200 && response.status < 300) {
                console.log("delete VMs successfully. HTTP code: %d. response: %s", response.status, text);
                setTimeout(unlockDeleteVM, 2000);
            } else {
                console.log("delete VMs failed. HTTP code: %d. response: %s", response.status, text);
                let deletingText = document.getElementById("textForDeleting");
                deletingText.textContent = `Deleting failed, HTTP code is ${response.status}, error is: ${text}`;
            }
        })
    }).catch(error => {
        console.error("Error:", error);
    });

}