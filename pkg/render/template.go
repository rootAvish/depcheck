package render

const HTMLTemplate = `<html>
<head>
   <link rel="stylesheet" href="https://stackpath.bootstrapcdn.com/bootstrap/4.2.1/css/bootstrap.min.css" integrity="sha384-GJzZqFGwb1QTTN6wy59ffF1BuGJpLSa9DkKMp0DgiMDm4iYMj70gZWKYbI706tWS" crossorigin="anonymous">
</head>
<body>
<h3>{{.Payload.Image}}</h3>
<p class="lead small text-right">
{{.Payload.Digest}}<br/>
Generated at {{.LastUpdate}}
</p>
<table class="table table-striped">
<tbody>
{{range $.Repositories}}
<tr>
<td><b>{{.Name}}</b></td></td>
<td>
<table class="table table-sm table-hover small">
<thead>
<tr class="text-left"><th>Dependency ({{.ManifestType}})</th><th>Version/Branch</th><th>Commit</th><th>Repository</th></tr>
</thead>
<tbody>
{{range .Dependencies}}
<tr>
<td>{{.Name}}</td><td>{{.Version}}</td><td><code>{{.Digest}}</code></td><td><a href="{{.Repository}}">{{.Repository}}</a></td>
</tr>
{{end}}
</tbody>
</table>
</td>
</tr>
{{end}}
</tbody>
</table>
</body>
</html>
`
