<!doctype html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport"
          content="width=device-width, user-scalable=no, initial-scale=1.0, maximum-scale=1.0, minimum-scale=1.0">
    <meta http-equiv="X-UA-Compatible" content="ie=edge">
    <title>Create Virtual Machine Success</title>

    <link rel="stylesheet" href="/static/css/style.css">

    <meta http-equiv="refresh" content="3; url=/cloud/{{.cloudName}}" />
</head>
<body>
{{template "/public/header.tpl" .}}
<h2>Create Virtual Machine Success</h2>
<h3>Redirect to the page of cloud: [{{.cloudName}}] after 3 seconds</h3>

</body>
</html>