<!doctype html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport"
          content="width=device-width, user-scalable=no, initial-scale=1.0, maximum-scale=1.0, minimum-scale=1.0">
    <meta http-equiv="X-UA-Compatible" content="ie=edge">
    <title>New Application</title>

    <link rel="stylesheet" href="/static/css/style.css">

    <!--We allow adding several nodes in one time-->
    <script src="/static/js/addK8sNodes.js"></script>

</head>
<body onload="addOneNode()">
    <a href="/k8sNode">
        <button>Back</button>
    </a>
    <br>
    <h2>New Nodes</h2>

    <form id="nodesInfo" action="/k8sNode/doAdd" method="post" onsubmit="whileAddingNodes()">

        <!--submit the container Number-->
        <input type="hidden" id="newNodeNum" name="newNodeNumber" value="0">
        <br>
        <b>Nodes to add:</b>
        <br>
        <br id="nodesStart">

        <button id="addOneButton" type="button" onclick="addOneNode()">One More</button>
        <button id="deleteOneButton" type="button" onclick="deleteOneNode('')">One Less</button>
        <br><br>

        <input id ="nodesInfoSubmit" type="submit" value="Submit">
    </form>

</body>
</html>