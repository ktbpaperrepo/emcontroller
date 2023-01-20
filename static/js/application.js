'use strict';

// For safety, we do not allow users to delete applications too frequently
let deleteAppLock = false;

function lockDeleteApp() {
    deleteAppLock = true;
}

function unlockDeleteApp() {
    deleteAppLock = false;
    location.reload() // after deleting, refresh the page
}

// original html does not support to send PUT or DELETE request
function deleteApp(appName, statusID) {
    if (deleteAppLock) {
        console.log("Another deleting is executing, please try again after a few seconds");
        return;
    }
    lockDeleteApp();
    let appStatus = document.getElementById(statusID);
    appStatus.innerText = "Deleting";
    let xmlhttp = new XMLHttpRequest();
    xmlhttp.open("DELETE", `/application/${appName}`);
    xmlhttp.send();
    console.log("delete %s request has been sent", appName);
    xmlhttp.onreadystatechange = function(){
        if(this.readyState==4 && this.status==200) {
            console.log("delete %s response: %s", appName, xmlhttp.responseText);
            // reserve 2s for deleting
            setTimeout(unlockDeleteApp, 2000);
        }
        console.log("onreadystatechange this.readyState: %O, this.status: %O", this.readyState, this.status);
    }
}

