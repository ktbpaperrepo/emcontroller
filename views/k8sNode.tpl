<!doctype html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport"
          content="width=device-width, user-scalable=no, initial-scale=1.0, maximum-scale=1.0, minimum-scale=1.0">
    <meta http-equiv="X-UA-Compatible" content="ie=edge">
    <title>Kubernetes Node</title>

    <link rel="stylesheet" href="/static/css/style.css">
    <link rel="stylesheet" href="/static/css/button.css">

    <script src="/static/js/k8sNode.js"></script>

</head>
<body>
    {{template "/public/header.tpl" .}}
    <h2>Kubernetes Node</h2>

    <br>
    <button class="button buttonFont change" onclick="window.open('/k8sNode/add', '_blank')"><h3>Add Nodes</h3></button>
    <br>

    <h3>Existing Kubernetes Nodes (Master Nodes not Shown)</h3>

    <button id="deleteSelectedButton" type="button" onclick="deleteBatchNodes()">Delete Selected Kubernetes Nodes</button>

    <table border = 1>
        <tr>
            <th rowspan="2"></th>
            <th rowspan="2"></th>
            <th rowspan="2">Name</th>
            <th rowspan="2">IP address</th>
            <th colspan="3">Resources (used/total)</th>
            <th rowspan="2">Status</th>
        </tr>
        <tr>
            <th>CPU Logical Core</th>
            <th>Memory (MB)</th>
            <th>Storage (GB)</th>
        </tr>
        {{range $nodeIdx, $node := .k8sNodeList}}
            {{$statusID := printf "nodeStatus%s" $node.Name}}
            <tr>
                <td><input type="checkbox" class="nodeCheckbox"></td>
                <td><button type="button" onclick="deleteNode('{{$node.Name}}', '{{$statusID}}')">Delete</button></td>
                <td>{{$node.Name}}</td>
                <td>{{$node.IP}}</td>
                <td>{{$node.UsedResources.CpuCore}}/{{$node.TotalResources.CpuCore}}</td>
                <td>{{$node.UsedResources.Memory}}/{{$node.TotalResources.Memory}}</td>
                <td>{{$node.UsedResources.Storage}}/{{$node.TotalResources.Storage}}</td>
                <td id="{{$statusID}}">{{$node.Status}}</td>
            </tr>
        {{end}}
    </table>

</body>
</html>