<!doctype html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport"
          content="width=device-width, user-scalable=no, initial-scale=1.0, maximum-scale=1.0, minimum-scale=1.0">
    <meta http-equiv="X-UA-Compatible" content="ie=edge">
    <title>New VMs</title>

    <link rel="stylesheet" href="/static/css/style.css">

    <!--We allow adding several nodes in one time-->
    <script src="/static/js/newVms.js"></script>

</head>
<body onload="addOneVm()">
    {{template "/public/nameRules.tpl" .}}
    <a href="/vm">
        <button>Back</button>
    </a>
    <br>
    <h2>New VMs</h2>

    <form id="vmsInfo" action="/vm/doNew" method="post" onsubmit="whileAddingVms()">

        <!--submit the container Number-->
        <input type="hidden" id="newVmNum" name="newVmNumber" value="0">
        <br>
        <b>VMs to create:</b>
        <br>
        <br id="vmsStart">

        <button id="addOneButton" type="button" onclick="addOneVm()">One More</button>
        <button id="deleteOneButton" type="button" onclick="deleteOneVm('')">One Less</button>
        <br><br>

        <input id ="vmsInfoSubmit" type="submit" value="Submit">
    </form>

</body>
</html>