<!doctype html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport"
          content="width=device-width, user-scalable=no, initial-scale=1.0, maximum-scale=1.0, minimum-scale=1.0">
    <meta http-equiv="X-UA-Compatible" content="ie=edge">
    <title>Network State</title>

    <link rel="stylesheet" href="/static/css/style.css">
</head>
<body>
    {{template "/public/header.tpl" .}}
    <h2>Network State</h2>

    {{ if not .NetTestFuncOn }}
        <br>
        <h3>{{ .netTestFuncOffMsg }}</h3>
    {{ else }}
        <br>
        <h3>Round-Trip Time (RTT) with unit millisecond (ms)</h3>
        <h4>The data is measured at intervals of {{.NetTestPeriodSec}} seconds.</h4>

        {{ $netStateLen := len .netState }}

        {{ $nsKeys := .netStateKeys }}
        {{ $ns := .netState }}

        <table border = 1>
            <tr>
                <th rowspan="2" colspan="2">RTT (ms)</th>
                <th colspan="{{$netStateLen}}">Target Cloud</th>
            </tr>
            <tr>
                {{ range $idx, $cloud := $nsKeys }}
                    <th>{{$cloud}}</th>
                {{end}}
            </tr>
            <tr>
                <th rowspan="{{$netStateLen}}">Source Cloud</th>
                {{ $firstSKey := index $nsKeys 0 }}
                <th>{{$firstSKey}}</th>
                {{ range $tIdx, $tKey := $nsKeys }}
                    {{with $thisNs := index $ns $firstSKey $tKey}}
                        <td>{{$thisNs.Rtt}}</td>
                    {{end}}
                {{end}}
            </tr>
            {{ range $sIdx, $sKey := $nsKeys }}
                {{ if eq $sIdx 0 }}
                    {{continue}}
                {{ end }}
                <tr>
                <th>{{$sKey}}</th>
                {{ range $tIdx, $tKey := $nsKeys }}
                    {{with $thisNs := index $ns $sKey $tKey}}
                    <td>{{$thisNs.Rtt}}</td>
                    {{end}}
                {{end}}
                </tr>
            {{end}}
        </table>

    {{ end }}

</body>
</html>