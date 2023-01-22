<!doctype html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport"
          content="width=device-width, user-scalable=no, initial-scale=1.0, maximum-scale=1.0, minimum-scale=1.0">
    <meta http-equiv="X-UA-Compatible" content="ie=edge">
    <title>Cloud</title>

    <link rel="stylesheet" href="/static/css/style.css">
    <script src="/static/js/cloud.js"></script>
</head>
<body>
    {{template "/public/header.tpl" .}}
    <a href="/cloud">
        <button>Back</button>
    </a>
    <br>
    <h2>Cloud Name: [{{.cloudInfo.Name}}]&ensp;&ensp;&ensp;&ensp;Cloud Type: [{{.cloudInfo.Type}}]</h2>

    <br>
    <h3>Resources (used/total)</h3>

    <table border = 1>
        <tr>
            <th>CPU Logical Core</th>
            <th>Memory (MB)</th>
            <th>Storage (GB)</th>
            <th>Virtual Machine</th>
            <th>Volume</th>
            <th>Network Port</th>
        </tr>
        <tr>
            <td>{{.cloudInfo.Resources.InUse.VCpu}}/{{.cloudInfo.Resources.Limit.VCpu}}</td>
            <td>{{.cloudInfo.Resources.InUse.Ram}}/{{.cloudInfo.Resources.Limit.Ram}}</td>
            <td>{{.cloudInfo.Resources.InUse.Storage}}/{{.cloudInfo.Resources.Limit.Storage}}</td>
            {{if lt .cloudInfo.Resources.Limit.Vm 0.0}}
                <td></td>
            {{else}}
                <td>{{.cloudInfo.Resources.InUse.Vm}}/{{.cloudInfo.Resources.Limit.Vm}}</td>
            {{end}}
            {{if lt .cloudInfo.Resources.Limit.Volume 0.0}}
                <td></td>
            {{else}}
                <td>{{.cloudInfo.Resources.InUse.Volume}}/{{.cloudInfo.Resources.Limit.Volume}}</td>
            {{end}}
            {{if lt .cloudInfo.Resources.Limit.Port 0.0}}
                <td></td>
            {{else}}
                <td>{{.cloudInfo.Resources.InUse.Port}}/{{.cloudInfo.Resources.Limit.Port}}</td>
            {{end}}
        </tr>
    </table>

    <br>
    <h3>Create a new Virtual Machine</h3>
    <form id="uploadForm" method="POST" action="/cloud/{{.cloudInfo.Name}}/vm" enctype="multipart/form-data" onsubmit="whileCreatingVM()">
        <table>
            <tr>
                <th>VM Name:</th> <td><input type="text" id="newVmName" name="newVmName"/></td>
            </tr>
            <tr>
                <th>CPU Logical Core:</th> <td><input type="text" id="newVmVCpu" name="newVmVCpu"/></td>
            </tr>
            <tr>
                <th>Memory (MB):</th> <td><input type="text" id="newVmRam" name="newVmRam"/></td>
            </tr>
            <tr>
                <th>Storage (GB):</th> <td><input type="text" id="newVmStorage" name="newVmStorage"/></td>
            </tr>
        </table>
        <input id="createNewVm" type="submit" value="Create">
    </form>

    <br>
    <h3>Virtual Machines</h3>
    <table border = 1>
        <tr> <th rowspan="2"></th> <th rowspan="2">Name</th> <th rowspan="2">ID</th> <th rowspan="2">IP Addresses</th> <th colspan="3">Resources</th> <th rowspan="2">Status</th> </tr>
        <tr> <th>CPU Logical Core</th> <th>Memory (MB)</th> <th>Storage (GB)</th> </tr>
        {{range $vmIdx, $vm := .vmList}}
            {{$statusID := printf "vmStatus%s" $vm.ID}}
            <tr>
                <td><button type="button" onclick="deleteVM('{{$vm.Cloud}}', '{{$vm.ID}}', '{{$statusID}}')">Delete</button></td>
                <td>{{$vm.Name}}</td>
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