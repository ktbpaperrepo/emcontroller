'use strict';

// deleting an image should restart the docker registry container, so 2 deleting cannot be executed at one time
let deleteRepoLock = false;

function lockDeleteRepo() {
    deleteRepoLock = true;
}

function unlockDeleteRepo() {
    deleteRepoLock = false;
    location.reload() // after deleting, refresh the page
}

// original html does not support to send PUT or DELETE request
function deleteRepo(repository) {
    if (deleteRepoLock) {
        console.log("Another deleting is executing, please try again after a few seconds");
        return;
    }
    lockDeleteRepo();
    let xmlhttp = new XMLHttpRequest();
    xmlhttp.open("DELETE", `/image/${repository}`);
    xmlhttp.send();
    console.log("delete %s request has been sent", repository);
    xmlhttp.onreadystatechange = function(){
        if(this.readyState==4 && this.status==200) {
            console.log("delete %s response: %s", repository, xmlhttp.responseText);
            // reserve 1s for deleting
            setTimeout(unlockDeleteRepo, 1000);
        }
        console.log("onreadystatechange this.readyState: %O, this.status: %O", this.readyState, this.status);
    }
}

// while uploading an image, users cannot operate the upload part of the web
// can put this function in the onsubmit="whileUploading()" of the form.
function whileUploading() {
    let imageName = document.getElementById("imageName");
    let imageTag = document.getElementById("imageTag");
    let imageFile = document.getElementById("imageFile");
    let upload = document.getElementById("upload");

    imageName.setAttribute("readonly", "readonly");
    imageTag.setAttribute("readonly", "readonly");
    imageFile.setAttribute("readonly", "readonly");
    upload.setAttribute("disabled", "disabled");
}
// This is a better way, because Javascript stuff should not be inline in the HTML code.
// See https://stackoverflow.com/questions/5691054/disable-submit-button-on-form-submit
function initImagePage() {
    $("form#uploadForm").submit(function (event) {
        $(this).find(':input[type=text]').prop("readonly", "readonly");
        $(this).find(':input[type=file]').prop("readonly", "readonly");
        $(this).find(':input[type=submit]').prop("disabled", "disabled");
        let uploadButton = document.getElementById("upload");
        uploadButton.insertAdjacentHTML('afterend',"<p>Uploading, please wait ...</p>");
    });
}

