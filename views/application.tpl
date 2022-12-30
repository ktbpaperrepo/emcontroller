<!doctype html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport"
          content="width=device-width, user-scalable=no, initial-scale=1.0, maximum-scale=1.0, minimum-scale=1.0">
    <meta http-equiv="X-UA-Compatible" content="ie=edge">
    <title>Application</title>

    <link rel="stylesheet" href="/static/css/style.css">

    <script src="/static/js/application.js"></script>
</head>
<body>
    {{template "/public/header.tpl" .}}
    <h2>Application</h2>

    <br>
    <h3>New Application</h3>
    <form method="get" action="/newApplication">
        <input type="submit" value="New">
    </form>



    <br>
    <h3>Existing Applications</h3>
    <ul>
        {{range $key, $app := .applicationList}}
        <li>
            {{$app}} <button type="button" onclick="deleteApp('{{$app}}')">Delete</button>
        </li>
        {{end}}

    </ul>

</body>
</html>