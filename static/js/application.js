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


// delete multiple Applications
function deleteBatchApps() {
    if (deleteAppLock) {
        console.log("Another deleting is executing, please try again after a few seconds");
        return;
    }
    lockDeleteApp();

    let deleteSelectedButton = document.getElementById("deleteSelectedButton");
    deleteSelectedButton.insertAdjacentHTML('afterend',"<p id=\"textForDeleting\">Deleting selected Applications, please wait ...</p>");

    let appNamesToDelete = [];
    let appCheckboxes = document.getElementsByClassName("appCheckbox");

    for (let i = 0; i < appCheckboxes.length; i++) {
        if (appCheckboxes[i].checked) {
            // get the row of the table
            let row = appCheckboxes[i].parentNode.parentNode;

            // set the "Status" column as "Deleting"
            row.cells[5].innerText = "Deleting";

            // get the needed information to delete an application
            let appName = row.cells[2].textContent;

            // make the json body for the request
            appNamesToDelete.push(appName);

        }
    }

    // send http request to delete applications
    let resp = fetch("/application",{
        method: "DELETE",
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify(appNamesToDelete)
    })

    // The return type of fetch is Promise and the Promise can only be accessed in its then.
    // I need to access its status code and text.
    // - The status code can only be accessed in resp.then().
    // - The text can only be accessed in resp.text().then().
    // Therefore, I use a 2-layer then() to both access them.
    resp.then(response => {
        response.text().then(text => {
            if (response.status >= 200 && response.status < 300) {
                console.log("delete applications successfully. HTTP code: %d. response: %s", response.status, text);
                setTimeout(unlockDeleteApp, 2000);
            } else {
                console.log("delete applications failed. HTTP code: %d. response: %s", response.status, text);
                let deletingText = document.getElementById("textForDeleting");
                deletingText.textContent = `Deleting failed, HTTP code is ${response.status}, error is: ${text}`;
            }
        })
    }).catch(error => {
        console.error("Error:", error);
    });

}