<!doctype html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport"
          content="width=device-width, user-scalable=no, initial-scale=1.0, maximum-scale=1.0, minimum-scale=1.0">
    <meta http-equiv="X-UA-Compatible" content="ie=edge">
    <title>Kubernetes Node</title>

    <link rel="stylesheet" href="/static/css/style.css">

    <script src="/static/js/k8sNode.js"></script>

    <style>
        .button {
            background-color: white;
            border: 2px solid black;
            color: black;
            padding: 2px 20px;
            text-align: center;
            text-decoration: none;
            display: inline-block;
            margin: 4px 2px;
            cursor: pointer;
        }
        .buttonFont {font-size: 16px;}
        .change {
            border-color: #211A52;
            color: #211A52;
            background-color: #C0DAEF;
        }

        .change:hover {
            background: #211A52;
            color: white;
        }
    </style>
</head>
<body>
    {{template "/public/header.tpl" .}}
    <h2>Kubernetes Node</h2>

    <br>
    <button class="button buttonFont change" onclick="window.location.href='/k8sNode/add'"><h3>Add Nodes</h3></button>
    <br>

    <h3>Existing Kubernetes Nodes (Master Nodes not Shown)</h3>

    <table border = 1>
        <tr> <th></th> <th>Name</th> <th>IP address</th> <th>Status</th> </tr>
        {{range $nodeIdx, $node := .k8sNodeList}}
            {{$statusID := printf "nodeStatus%s" $node.Name}}
            <tr>
                <td><button type="button" onclick="deleteNode('{{$node.Name}}', '{{$statusID}}')">Delete</button></td>
                <td>{{$node.Name}}</td>
                <td>{{$node.IP}}</td>
                <td id="{{$statusID}}">{{$node.Status}}</td>
            </tr>
        {{end}}
    </table>

</body>
</html>