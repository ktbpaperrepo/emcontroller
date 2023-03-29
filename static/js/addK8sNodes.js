'use strict';

// We allow adding several nodes at one time
let nodeIndex = 0;

function generateNodeHTML() {
    return `
<div id="node${nodeIndex}">
    Name: <input type="text" name="node${nodeIndex}Name">&nbsp;&nbsp;&nbsp;IP Address: <input type="text" name="node${nodeIndex}IP">
</div>
`;
}

function addOneNode() {
    let addOneButton = document.getElementById("addOneButton");
    let nodeHTML = generateNodeHTML();
    addOneButton.insertAdjacentHTML('beforebegin', nodeHTML);

    let deleteOneButton = document.getElementById("deleteOneButton");
    deleteOneButton.setAttribute('onclick', `deleteOneNode('node${nodeIndex}')`);
    nodeIndex++;

    // update the new node Number in a form input for submission
    let newNodeNum = document.getElementById("newNodeNum");
    newNodeNum.setAttribute('value', String(nodeIndex));
}

function deleteOneNode(nodeElementID) {
    let deleteResult = deleteElement(nodeElementID);
    if (!deleteResult) {
        return;
    }
    nodeIndex--;
    let deleteOneButton = document.getElementById("deleteOneButton");
    deleteOneButton.setAttribute('onclick', `deleteOneNode('node${nodeIndex-1}')`);

    // update the new node Number in a form input for submission
    let newNodeNum = document.getElementById("newNodeNum");
    newNodeNum.setAttribute('value', String(nodeIndex));
}

function deleteElement(elementID) {
    let toBeDelete = document.getElementById(elementID);
    if (toBeDelete == null) {
        return false;
    }
    toBeDelete.remove();
    return true;
}

function whileAddingNodes() {
    let submitButton = document.getElementById("nodesInfoSubmit");

    submitButton.setAttribute("disabled", "disabled");
    submitButton.insertAdjacentHTML('afterend',`<p>Adding nodes, please wait ...</p>`);
}