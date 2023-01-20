<!doctype html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport"
          content="width=device-width, user-scalable=no, initial-scale=1.0, maximum-scale=1.0, minimum-scale=1.0">
    <meta http-equiv="X-UA-Compatible" content="ie=edge">
    <title>Image</title>

    <link rel="stylesheet" href="/static/css/style.css">
<!--jquery should before image.js in which jquery is used-->
    <script src="/static/js/jquery-3.6.3.js"></script>
    <script src="/static/js/image.js"></script>

</head>
<body onload="initImagePage()">
    {{template "/public/header.tpl" .}}
    <h2>Container Image</h2>

    <br>
    <h3>Upload a new Image</h3>
    <form id="uploadForm" method="POST" action="/upload" enctype="multipart/form-data">
        Set the RepoTag:<br>
        {{config "String" "dockerRegistryIP" ""}}:{{config "String" "dockerRegistryPort" ""}}/<input type="text" id="imageName" name="imageName">:<input type="text" id="imageTag" name="imageTag"><br>
        <input id="imageFile" name="imageFile" type="file"/>
        <input id="upload" type="submit" value="Upload">
    </form>

    <br>
    <h3>Existing Images</h3>
    {{$dockerRepo := .dockerRegistry}}
    <table border = 1>
        <tr> <th></th> <th>Repositories</th> <th>RepoTags</th> <th>Messages</th> </tr>
        {{range $repo, $tags := .imageList}}
            {{$messageID := printf "imageMessage%s" $repo}}
            <tr>
                <td><button type="button" onclick="deleteRepo('{{$repo}}', '{{$messageID}}')">Delete</button></td>
                <td>{{$repo}}</td>
                <td>
                    {{range $idx, $tag := $tags}}
                        {{$dockerRepo}}/{{$repo}}:{{$tag}} <br>
                    {{end}}
                </td>
                <td id="{{$messageID}}"></td>
            </tr>
        {{end}}
    </table>

</body>
</html>