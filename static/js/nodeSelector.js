'use strict';

let nodeSelectorIndex = 0;

function generateNodeSelectorHTML() {
    return `
<div id="nodeSelector${nodeSelectorIndex}">
    Key: <input type="text" name="nodeSelector${nodeSelectorIndex}Key">
    &nbsp;&nbsp;&nbsp;Value: <input type="text" name="nodeSelector${nodeSelectorIndex}Value">
</div>
`;
}

function addNodeSelector() {
    let addNodeSelectorButton = document.getElementById("addNodeSelectorButton");
    let nodeSelectorHTML = generateNodeSelectorHTML();
    addNodeSelectorButton.insertAdjacentHTML('beforebegin', nodeSelectorHTML);

    let deleteNodeSelectorButton = document.getElementById("deleteNodeSelectorButton");
    deleteNodeSelectorButton.setAttribute('onclick', `deleteNodeSelector('nodeSelector${nodeSelectorIndex}')`);
    nodeSelectorIndex++;

    // update the node selector Number
    let nodeSelectorNum = document.getElementById("nodeSelectorNum");
    nodeSelectorNum.setAttribute('value', String(nodeSelectorIndex));
}

function deleteNodeSelector(nodeSelectorElementID) {
    let deleteResult = deleteElement(nodeSelectorElementID);
    if (!deleteResult) {
        return;
    }
    nodeSelectorIndex--;
    let deleteNodeSelectorButton = document.getElementById("deleteNodeSelectorButton");
    deleteNodeSelectorButton.setAttribute('onclick', `deleteNodeSelector('nodeSelector${nodeSelectorIndex-1}')`);

    // update the node selector Number
    let nodeSelectorNum = document.getElementById("nodeSelectorNum");
    nodeSelectorNum.setAttribute('value', String(nodeSelectorIndex));
}

function deleteElement(elementID) {
    let toBeDelete = document.getElementById(elementID);
    if (toBeDelete == null) {
        return false;
    }
    toBeDelete.remove();
    return true;
}