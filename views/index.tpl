<!doctype html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport"
          content="width=device-width, user-scalable=no, initial-scale=1.0, maximum-scale=1.0, minimum-scale=1.0">
    <meta http-equiv="X-UA-Compatible" content="ie=edge">
    <title>{{.Website}}</title>
</head>
<body>
    {{template "/public/header.html" .}}
    <h2>1. bind basic data types in the template: string, num, bool</h2>

    <p>{{.title}}</p>
    <p>{{.num}}</p>
    <p>{{.flag}}</p>

    <br>

    <h2>2. bind the data of structures in the template</h2>
    <p>{{.article.Title}}</p>
    <p>{{.article.Content}}</p>
    <br>

    <h2>3. customize variables in the template</h2>
    {{$xxx := .title}}
    <p>{{$xxx}}</p>

    <h2>4. loop in template, range, loop slice</h2>
    <ul>
        {{range $key,$val := .sliceList}}
        <li>{{$key}}-----{{$val}}</li>
        {{end}}
    </ul>

    <h2>5. loop in template, range, loop Map</h2>
    <ul>
        {{range $key,$val := .userinfo}}
        <li>{{$key}}-----{{$val}}</li>
        {{end}}
    </ul>

    <h2>6. loop in template, loop slice with the type of structure</h2>
    <ul>
        {{range $key, $val := .articleList}}
            <li>
                {{$key}}---{{$val.Title}}---{{$val.Content}}
            </li>
        {{end}}
    </ul>

    <h2>7. loop in template, another definition method of slice with the type of structure</h2>
    <ul>
        {{range $key, $val := .cmdList}}
        <li>
            {{$key}}---{{$val.Title}}
        </li>
        {{end}}
    </ul>

    <h2>8. conditions in template</h2>
    {{if .isLogin}}
        <p>isLogin is true</p>
    {{end}}

    {{if .isHome}}
        <p>isHome is true</p>
    {{else}}
        <p>isHome is false</p>
    {{end}}

    {{if .isHome}}
    <p>isHome is true</p>
    {{else if .isAbout}}
    <p>isAbout is true</p>
    {{end}}

    {{if .isHome}}
    <p>isHome is true</p>
    {{else}}
        {{if .isAbout}}
            <p>isAbout is true 1111</p>
        {{end}}
    {{end}}

    <h2>9. if condition eq / ne / lt / le / gt / ge</h2>
    {{if gt .n1 .n2}}
        <p>n1 is larger than n2</p>
    {{end}}

    {{if eq .n1 .n2}}
        <p>n1 is equal to n2</p>
    {{else}}
        <p>n1 is not equal to n2</p>
    {{end}}

    {{if ne .n1 .n2}}
        <p>n1 is not equal to n2 ----</p>
    {{else}}
        <p>n1 is equal to n2 ----</p>
    {{end}}


    <h2>10. define template</h2>
    {{define "aaa"}}
        <h4>This is a defined code block</h4>
        <p>111</p>
        <p>2222211</p>
    {{end}}

    <div>
        {{template "aaa" .}}
    </div>
    <hr>
    <div>
        {{template "aaa" .}}
    </div>


    <h2>11. define template outside</h2>
    {{template "/public/footer.html" .}}

    <h2>// 12. use self-defined template functions</h2>
    <p>{{.unix | unixToDate}}</p>

</body>
</html>