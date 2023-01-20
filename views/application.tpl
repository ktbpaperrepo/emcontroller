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
        <input type="radio" name="mode" value="basic" checked="checked" />Basic Mode
        <input type="radio" name="mode" value="advanced" />Advanced Mode <br>
        <input type="submit" value="New">
    </form>



    <br>
    <h3>Existing Applications</h3>

    <table border = 1>
        <tr> <th></th> <th>App Name</th> <th>Internal Access</th> <th>External Access</th> <th>Status</th> </tr>
        {{range $appIdx, $app := .applicationList}}
            {{$statusID := printf "appStatus%s" $app.AppName}}
            <tr>
                <td><button type="button" onclick="deleteApp('{{$app.AppName}}', '{{$statusID}}')">Delete</button></td>
                <td>{{$app.AppName}}</td>
                <td>
                    {{if not (eq $app.ClusterIP "" "None") }}
                        {{range $idx, $svcPort := $app.SvcPort}}
                        {{$app.SvcName}}:{{$svcPort}} <br>
                        {{$app.ClusterIP}}:{{$svcPort}} <br>
                        {{end}}
                    {{end}}
                </td>
                <td>
                    {{if ne $app.NodePortIP ""}}
                        {{range $idx, $nodePort := $app.NodePort}}
                        {{$app.NodePortIP}}:{{$nodePort}} <br>
                        {{end}}
                    {{end}}
                </td>
                <td id="{{$statusID}}">{{$app.Status}}</td>
            </tr>
        {{end}}
    </table>

</body>
</html>