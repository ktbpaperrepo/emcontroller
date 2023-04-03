'use strict';

// We allow creating several vms at one time
let vmIndex = 0;

function generateVmHTML() {
    return `
<div id="vm${vmIndex}">
    Name: <input type="text" name="vm${vmIndex}Name">
    &nbsp;&nbsp;&nbsp;Cloud Name: <input type="text" name="vm${vmIndex}CloudName">
    &nbsp;&nbsp;&nbsp;CPU Logical Core: <input type="text" name="vm${vmIndex}VCpu">
    &nbsp;&nbsp;&nbsp;Memory (MB): <input type="text" name="vm${vmIndex}Ram">
    &nbsp;&nbsp;&nbsp;Storage (GB): <input type="text" name="vm${vmIndex}Storage">
</div>
`;
}

function addOneVm() {
    let addOneButton = document.getElementById("addOneButton");
    let vmHTML = generateVmHTML();
    addOneButton.insertAdjacentHTML('beforebegin', vmHTML);

    let deleteOneButton = document.getElementById("deleteOneButton");
    deleteOneButton.setAttribute('onclick', `deleteOneVm('vm${vmIndex}')`);
    vmIndex++;

    // update the new VM Number in a form input for submission
    let newVmNum = document.getElementById("newVmNum");
    newVmNum.setAttribute('value', String(vmIndex));
}

function deleteOneVm(vmElementID) {
    let deleteResult = deleteElement(vmElementID);
    if (!deleteResult) {
        return;
    }
    vmIndex--;
    let deleteOneButton = document.getElementById("deleteOneButton");
    deleteOneButton.setAttribute('onclick', `deleteOneVm('vm${vmIndex-1}')`);

    // update the new VM Number in a form input for submission
    let newVmNum = document.getElementById("newVmNum");
    newVmNum.setAttribute('value', String(vmIndex));
}

function deleteElement(elementID) {
    let toBeDelete = document.getElementById(elementID);
    if (toBeDelete == null) {
        return false;
    }
    toBeDelete.remove();
    return true;
}

function whileAddingVms() {
    let submitButton = document.getElementById("vmsInfoSubmit");

    submitButton.setAttribute("disabled", "disabled");
    submitButton.insertAdjacentHTML('afterend',`<p>Creating Virtual Machines, please wait ...</p>`);
}