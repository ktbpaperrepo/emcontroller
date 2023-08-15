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


// delete multiple Kubernetes nodes
function deleteBatchNodes() {
    if (deleteNodeLock) {
        console.log("Another deleting is executing, please try again after a few seconds");
        return;
    }
    lockDeleteNode();

    let deleteSelectedButton = document.getElementById("deleteSelectedButton");
    deleteSelectedButton.insertAdjacentHTML('afterend',"<p id=\"textForDeleting\">Deleting selected Kubernetes Nodes, please wait ...</p>");

    let nodeNamesToDelete = [];
    let nodeCheckboxes = document.getElementsByClassName("nodeCheckbox");

    for (let i = 0; i < nodeCheckboxes.length; i++) {
        if (nodeCheckboxes[i].checked) {
            // get the row of the table
            let row = nodeCheckboxes[i].parentNode.parentNode;

            // set the "Status" column as "Deleting"
            row.cells[7].innerText = "Deleting";

            // get the needed information to delete a Kubernetes node
            let nodeName = row.cells[2].textContent;

            // make the json body for the request
            nodeNamesToDelete.push(nodeName);

        }
    }

    // send http request to delete Kubernetes nodes
    let resp = fetch("/k8sNode",{
        method: "DELETE",
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify(nodeNamesToDelete)
    })

    // The return type of fetch is Promise and the Promise can only be accessed in its then.
    // I need to access its status code and text.
    // - The status code can only be accessed in resp.then().
    // - The text can only be accessed in resp.text().then().
    // Therefore, I use a 2-layer then() to both access them.
    resp.then(response => {
        response.text().then(text => {
            if (response.status >= 200 && response.status < 300) {
                console.log("delete Kubernetes nodes successfully. HTTP code: %d. response: %s", response.status, text);
                setTimeout(unlockDeleteNode, 2000);
            } else {
                console.log("delete Kubernetes nodes failed. HTTP code: %d. response: %s", response.status, text);
                let deletingText = document.getElementById("textForDeleting");
                deletingText.textContent = `Deleting failed, HTTP code is ${response.status}, error is: ${text}`;
            }
        })
    }).catch(error => {
        console.error("Error:", error);
    });

}