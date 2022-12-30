'use strict';

// original html does not support to send PUT or DELETE request
function deleteApp(appName) {
    let xmlhttp = new XMLHttpRequest()
    xmlhttp.open("DELETE", `/application/${appName}`)
    xmlhttp.send()
    xmlhttp.onreadystatechange = function(){
        if(this.readyState==4 && this.status==200) {
            console.log(xmlhttp.responseText)
        }
    }
}

