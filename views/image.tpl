<!doctype html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport"
          content="width=device-width, user-scalable=no, initial-scale=1.0, maximum-scale=1.0, minimum-scale=1.0">
    <meta http-equiv="X-UA-Compatible" content="ie=edge">
    <title>Image</title>

    <link rel="stylesheet" href="/static/css/style.css">

    <script src="/static/js/image.js"></script>

</head>
<body>
    {{template "/public/header.tpl" .}}
    <h2>Image</h2>

    <br>
    <h3>Upload a new Image</h3>
    <form method="POST" action="/upload" enctype="multipart/form-data">
        Set the RepoTag:<br>
        {{config "String" "dockerRegistryIP" ""}}:{{config "String" "dockerRegistryPort" ""}}/<input type="text" name="imageName">:<input type="text" name="imageTag"><br>
        <input name="imageFile" type="file"/>
        <input type="submit" value="Upload">
    </form>

    <br>
    <h3>Existing Images</h3>
    {{$dockerRepo := .dockerRegistry}}
    <table border = 1>
        <tr> <th></th> <th>Repositories</th> <th>RepoTags</th> </tr>
        {{range $repo, $tags := .imageList}}
            <tr>
                <td><button type="button" onclick="deleteRepo('{{$repo}}')">Delete</button></td>
                <td>{{$repo}}</td>
                <td>
                    {{range $idx, $tag := $tags}}
                        {{$dockerRepo}}/{{$repo}}:{{$tag}} <br>
                    {{end}}
                </td>
            </tr>
        {{end}}
    </table>

</body>
</html>