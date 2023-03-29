'use strict';

// For safety, we do not allow users to delete nodes too frequently
let deleteNodeLock = false;

function lockDeleteNode() {
    deleteNodeLock = true;
}

function unlockDeleteNode() {
    deleteNodeLock = false;
    location.reload() // after deleting, refresh the page
}

// original html does not support to send PUT or DELETE request
function deleteNode(nodeName, statusID) {
    if (deleteNodeLock) {
        console.log("Another deleting is executing, please try again after a few seconds");
        return;
    }
    lockDeleteNode();
    let nodeStatus = document.getElementById(statusID);
    nodeStatus.innerText = "Deleting";
    let xmlhttp = new XMLHttpRequest();
    xmlhttp.open("DELETE", `/k8sNode/${nodeName}`);
    xmlhttp.send();
    console.log("delete %s request has been sent", nodeName);
    xmlhttp.onreadystatechange = function(){
        if (this.readyState==4) {
            if (this.status==200) {
                console.log("delete %s response: %s", nodeName, xmlhttp.responseText);
                // reserve 2s for deleting
                setTimeout(unlockDeleteNode, 2000);
            } else {
                nodeStatus.innerText = `Delete Error:\r\nreadyState: ${this.readyState}\r\nstatus: ${this.status}`;
            }
        }
        console.log("onreadystatechange this.readyState: %O, this.status: %O", this.readyState, this.status);
    }
}

