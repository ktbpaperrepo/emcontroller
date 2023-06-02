<!doctype html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport"
          content="width=device-width, user-scalable=no, initial-scale=1.0, maximum-scale=1.0, minimum-scale=1.0">
    <meta http-equiv="X-UA-Compatible" content="ie=edge">
    <title>{{.Website}}</title>

    <link rel="stylesheet" href="/static/css/style.css">

</head>
<body>
    {{template "/public/header.tpl" .}}
    <h2 style="text-align:center">{{.Website}}</h2>
    <br>
    <h3 style="text-align:center">{{.VersionInfo}}</h3>
    <br>
    <br>
    <h2 style="text-align:center">Please use your browser's <span style="color: red; font-weight: bold;">Incognito Mode</span> to visit this website; otherwise, some cached resources will make some functions abnormal. </h2>

</body>
</html>