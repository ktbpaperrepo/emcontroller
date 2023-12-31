<!doctype html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport"
          content="width=device-width, user-scalable=no, initial-scale=1.0, maximum-scale=1.0, minimum-scale=1.0">
    <meta http-equiv="X-UA-Compatible" content="ie=edge">
    <title>Virtual Machine</title>

    <link rel="stylesheet" href="/static/css/style.css">
    <link rel="stylesheet" href="/static/css/button.css">

    <script src="/static/js/vm.js"></script>
</head>
<body>
    {{template "/public/header.tpl" .}}
    <h2>Virtual Machine</h2>

    <br>
    <button class="button buttonFont change" onclick="window.open('/vm/new', '_blank')"><h3>Create VMs</h3></button>
    <br>

    <br>
    <h3>Existing VMs in All Clouds</h3>

    <button id="deleteSelectedButton" type="button" onclick="deleteBatchVMs()">Delete Selected VMs</button>

    <table border = 1>
        <tr>
            <th rowspan="2"></th>
            <th rowspan="2"></th>
            <th rowspan="2">Multi-Cloud<br>Manager Create</th>
            <th rowspan="2">Name</th>
            <th rowspan="2">Cloud Type</th>
            <th rowspan="2">Cloud</th>
            <th rowspan="2">ID</th>
            <th rowspan="2">IP Addresses</th>
            <th colspan="3">Resources</th>
            <th rowspan="2">Status</th>
        </tr>
        <tr>
            <th>CPU Logical Core</th>
            <th>Memory (MB)</th>
            <th>Storage (GB)</th>
        </tr>

        {{range $vmIdx, $vm := .allVms}}
            {{$statusID := printf "vmStatus-%s-%s" $vm.Cloud $vm.ID}}
            <tr>
                <td><input type="checkbox" class="vmCheckbox"></td>
                <td><button type="button" onclick="deleteVM('{{$vm.Cloud}}', '{{$vm.ID}}', '{{$statusID}}')">Delete</button></td>
                {{$colorMcmCreated := "red"}}
                {{if $vm.McmCreate}}
                    {{$colorMcmCreated = "green"}}
                {{end}}
                <td style="color:{{$colorMcmCreated}}">{{$vm.McmCreate}}</td>
                <td>{{$vm.Name}}</td>
                <td>{{$vm.CloudType}}</td>
                <td>{{$vm.Cloud}}</td>
                <td>{{$vm.ID}}</td>
                <td>
                    {{range $idx, $ip := $vm.IPs}}
                        {{$ip}} <br>
                    {{end}}
                </td>
                <td>{{$vm.VCpu}}</td>
                <td>{{$vm.Ram}}</td>
                <td>{{$vm.Storage}}</td>
                <td id="{{$statusID}}">{{$vm.Status}}</td>
            </tr>
        {{end}}

    </table>

</body>
</html>