<!doctype html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport"
          content="width=device-width, user-scalable=no, initial-scale=1.0, maximum-scale=1.0, minimum-scale=1.0">
    <meta http-equiv="X-UA-Compatible" content="ie=edge">
    <title>Cloud</title>

    <link rel="stylesheet" href="/static/css/style.css">
</head>
<body>
    {{template "/public/header.tpl" .}}
    <h2>Cloud</h2>

    <br>
    <h3>All Clouds</h3>

    <table border = 1>
        <tr>
            <th rowspan="2">Name</th>
            <th rowspan="2">Type</th>
            <th colspan="6">Resources (used/total)</th>
        </tr>
        <tr>
<!--            <th></th> -->
<!--            <th></th> -->
            <th>CPU Logical Core</th>
            <th>Memory (MB)</th>
            <th>Storage (GB)</th>
            <th>Virtual Machine</th>
            <th>Volume</th>
            <th>Network Port</th>
        </tr>
        {{range $cloudIdx, $cloud := .cloudList}}
            <tr>
                <td><a href="/cloud/{{$cloud.Name}}">{{$cloud.Name}}</a></td>
                <td>{{$cloud.Type}}</td>
                <td>{{$cloud.Resources.InUse.VCpu}}/{{$cloud.Resources.Limit.VCpu}}</td>
                <td>{{$cloud.Resources.InUse.Ram}}/{{$cloud.Resources.Limit.Ram}}</td>
                <td>{{$cloud.Resources.InUse.Storage}}/{{$cloud.Resources.Limit.Storage}}</td>
                {{if lt $cloud.Resources.Limit.Vm 0.0}}
                    <td></td>
                {{else}}
                    <td>{{$cloud.Resources.InUse.Vm}}/{{$cloud.Resources.Limit.Vm}}</td>
                {{end}}
                {{if lt $cloud.Resources.Limit.Volume 0.0}}
                    <td></td>
                {{else}}
                    <td>{{$cloud.Resources.InUse.Volume}}/{{$cloud.Resources.Limit.Volume}}</td>
                {{end}}
                {{if lt $cloud.Resources.Limit.Port 0.0}}
                    <td></td>
                {{else}}
                    <td>{{$cloud.Resources.InUse.Port}}/{{$cloud.Resources.Limit.Port}}</td>
                {{end}}
            </tr>
        {{end}}
    </table>

</body>
</html>