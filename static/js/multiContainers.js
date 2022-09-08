'use strict';

// There may be several containers in a pod
let containerIndex = 0;

function generateContainerHTML() {
    return `
<div id="container${containerIndex}">
    <p>Container ${containerIndex+1}</p>
    Container Name: <input type="text" name="container${containerIndex}Name"> <br><br>
    Image RepoTag: <input type="text" name="container${containerIndex}Image"> <br><br>
    Resources requests:
    <ul>
        <li>Memory: <input type="text" name="container${containerIndex}RequestMemory"></li>
        <li>CPU: <input type="text" name="container${containerIndex}RequestCPU"></li>
        <li>Ephemeral Storag: <input type="text" name="container${containerIndex}RequestEphemeralStorage"></li>
    </ul>
    Resources Limits:
    <ul>
        <li>Memory: <input type="text" name="container${containerIndex}LimitMemory"></li>
        <li>CPU: <input type="text" name="container${containerIndex}LimitCPU"></li>
        <li>Ephemeral Storag: <input type="text" name="container${containerIndex}LimitEphemeralStorage"></li>
    </ul>
    Commands:
    <!--submit the Command Number-->
    <input type="hidden" id="container${containerIndex}CommandNum" name="container${containerIndex}CommandNumber" value="0">
    <ul id="container${containerIndex}Commands">
        <button id="container${containerIndex}AddCommandButton" type="button" onclick="addCommand('container${containerIndex}')">Add</button>
        <button id="container${containerIndex}DeleteCommandButton" type="button" onclick="deleteCommand('container${containerIndex}')">Delete</button>
    </ul>
    Args:
    <!--submit the Arg Number-->
    <input type="hidden" id="container${containerIndex}ArgNum" name="container${containerIndex}ArgNumber" value="0">
    <ul id="container${containerIndex}Args">
        <button id="container${containerIndex}AddArgButton" type="button" onclick="addArg('container${containerIndex}')">Add</button>
        <button id="container${containerIndex}DeleteArgButton" type="button" onclick="deleteArg('container${containerIndex}')">Delete</button>
    </ul>
    Ports:
    <!--submit the Port Number-->
    <input type="hidden" id="container${containerIndex}PortNum" name="container${containerIndex}PortNumber" value="0">
    <ul id="container${containerIndex}Ports">
        <button id="container${containerIndex}AddPortButton" type="button" onclick="addPort('container${containerIndex}')">Add</button>
        <button id="container${containerIndex}DeletePortButton" type="button" onclick="deletePort('container${containerIndex}')">Delete</button>
    </ul>

    <br>
    <br>
</div>
`;
}

function generateCommandHTML(containerElementID) {
    let container = document.getElementById(containerElementID);
    return `
        <li id="${containerElementID}Command${container.commandIndex}">
            <input type="text" name="${containerElementID}Command${container.commandIndex}">
        </li>
    `;
}

function generateArgHTML(containerElementID) {
    let container = document.getElementById(containerElementID);
    return `
        <li id="${containerElementID}Arg${container.argIndex}">
            <input type="text" name="${containerElementID}Arg${container.argIndex}">
        </li>
    `;
}

function generatePortHTML(containerElementID) {
    let container = document.getElementById(containerElementID);
    return `
        <li id="${containerElementID}Port${container.portIndex}">
            ContainerPort: <input type="text" name="${containerElementID}Port${container.portIndex}ContainerPort"><br>
            Name: <input type="text" name="${containerElementID}Port${container.portIndex}Name"><br>
            Protocol: <input type="text" name="${containerElementID}Port${container.portIndex}Protocol"><br>
            ServicePort: <input type="text" name="${containerElementID}Port${container.portIndex}ServicePort"><br>
            NodePort: <input type="text" name="${containerElementID}Port${container.portIndex}NodePort"><br>
        </li>
        <br>
    `;
}

function addContainer() {
    let addContainerButton = document.getElementById("addContainerButton");
    let containerHTML = generateContainerHTML()
    addContainerButton.insertAdjacentHTML('beforebegin', containerHTML);

    // initialize
    let newContainer = document.getElementById(`container${containerIndex}`);
    newContainer.commandIndex = 0;
    newContainer.argIndex = 0;
    newContainer.portIndex = 0;

    let deleteContainerButton = document.getElementById("deleteContainerButton");
    deleteContainerButton.setAttribute('onclick', `deleteContainer('container${containerIndex}')`);
    containerIndex++;

    // update the container Number in a form input for submission
    let containerNumber = document.getElementById("containerNum");
    containerNumber.setAttribute('value', String(containerIndex));
}

function deleteContainer(containerElementID) {
    let deleteResult = deleteElement(containerElementID);
    if (!deleteResult) {
        return;
    }
    containerIndex--;
    let deleteContainerButton = document.getElementById("deleteContainerButton");
    deleteContainerButton.setAttribute('onclick', `deleteContainer('container${containerIndex-1}')`);

    // update the container Number in a form input for submission
    let containerNum = document.getElementById("containerNum");
    containerNum.setAttribute('value', String(containerIndex));
}

function addCommand(containerElementID) {
    let container = document.getElementById(containerElementID);
    let commandHTML = generateCommandHTML(containerElementID);
    let addCommandButton = document.getElementById(`${containerElementID}AddCommandButton`);
    addCommandButton.insertAdjacentHTML('beforebegin', commandHTML);
    container.commandIndex++;

    // update the command Number of this container in a form input for submission
    let commandNum = document.getElementById(`${containerElementID}CommandNum`);
    commandNum.setAttribute('value', String(container.commandIndex));
}

function deleteCommand(containerElementID) {
    let container = document.getElementById(containerElementID);
    let lastCommandElementID = container.commandIndex - 1;
    if (lastCommandElementID < 0) {
        return;
    }
    let lastCommand = document.getElementById(`${containerElementID}Command${lastCommandElementID}`);
    lastCommand.remove();
    container.commandIndex--;

    // update the command Number of this container in a form input for submission
    let commandNum = document.getElementById(`${containerElementID}CommandNum`);
    commandNum.setAttribute('value', String(container.commandIndex));
}

function addArg(containerElementID) {
    let container = document.getElementById(containerElementID);
    let argHTML = generateArgHTML(containerElementID);
    let addArgButton = document.getElementById(`${containerElementID}AddArgButton`);
    addArgButton.insertAdjacentHTML('beforebegin', argHTML);
    container.argIndex++;

    // update the arg Number of this container in a form input for submission
    let argNum = document.getElementById(`${containerElementID}ArgNum`);
    argNum.setAttribute('value', String(container.argIndex));
}

function deleteArg(containerElementID) {
    let container = document.getElementById(containerElementID);
    let lastArgElementID = container.argIndex - 1;
    if (lastArgElementID < 0) {
        return;
    }
    let lastArg = document.getElementById(`${containerElementID}Arg${lastArgElementID}`);
    lastArg.remove();
    container.argIndex--;

    // update the arg Number of this container in a form input for submission
    let argNum = document.getElementById(`${containerElementID}ArgNum`);
    argNum.setAttribute('value', String(container.argIndex));
}

function addPort(containerElementID) {
    let container = document.getElementById(containerElementID)
    let portHTML = generatePortHTML(containerElementID)
    let addPortButton = document.getElementById(`${containerElementID}AddPortButton`)
    addPortButton.insertAdjacentHTML('beforebegin', portHTML)
    container.portIndex++;

    // update the port Number of this container in a form input for submission
    let portNum = document.getElementById(`${containerElementID}PortNum`);
    portNum.setAttribute('value', String(container.portIndex));
}

function deletePort(containerElementID) {
    let container = document.getElementById(containerElementID);
    let lastPortElementID = container.portIndex - 1;
    if (lastPortElementID < 0) {
        return;
    }
    let lastPort = document.getElementById(`${containerElementID}Port${lastPortElementID}`);
    lastPort.remove();
    container.portIndex--;

    // update the port Number of this container in a form input for submission
    let portNum = document.getElementById(`${containerElementID}PortNum`);
    portNum.setAttribute('value', String(container.portIndex));
}

function deleteElement(elementID) {
    let toBeDelete = document.getElementById(elementID);
    if (toBeDelete == null) {
        return false;
    }
    toBeDelete.remove();
    return true;
}

