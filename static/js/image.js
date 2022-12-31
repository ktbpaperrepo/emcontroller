'use strict';

// original html does not support to send PUT or DELETE request
function deleteRepo(repository) {
    let xmlhttp = new XMLHttpRequest()
    xmlhttp.open("DELETE", `/image/${repository}`)
    xmlhttp.send()
    xmlhttp.onreadystatechange = function(){
        if(this.readyState==4 && this.status==200) {
            console.log(xmlhttp.responseText)
        }
    }
}

