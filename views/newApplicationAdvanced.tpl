<!doctype html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport"
          content="width=device-width, user-scalable=no, initial-scale=1.0, maximum-scale=1.0, minimum-scale=1.0">
    <meta http-equiv="X-UA-Compatible" content="ie=edge">
    <title>New Application</title>

    <link rel="stylesheet" href="/static/css/style.css">

    <!--There may be several containers in a pod-->
    <script src="/static/js/multiContainers.js"></script>

</head>
<body onload="initAdvanced()">
    <a href="/application">
        <button>Back</button>
    </a>
    <br>
    <h2>New Application</h2>

    <form id="appInfo" action="/doNewApplication" method="post">
        Name: <input type="text" name="name"> <br><br>
        Replicas: <input type="text" name="replicas"> <br><br>

        <!--submit the container Number-->
        <input type="hidden" id="containerNum" name="containerNumber" value="0">
        <br>
        <b>Containers in each replica:</b>
        <br>
        <br id="containerStart">

        <button id="addContainerButton" type="button" onclick="addContainer()">Add Container</button>
        <button id="deleteContainerButton" type="button" onclick="deleteContainer('')">Delete Container</button>
        <br><br>

        <input id ="appInfoSubmit" type="submit" value="Create">
    </form>


</body>
</html>