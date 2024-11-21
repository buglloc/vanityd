package httpd

import "html/template"

var indexTmpl = template.Must(
	template.
		New("index").
		Parse(
			`
<html lang="en">
<head>
    <meta charset="utf-8">
    <title>@buglloc projects</title>
	<link rel="shortcut icon" href="/static/favicon.ico">
    <link rel="stylesheet" href="https://fonts.googleapis.com/css?family=Open+Sans:300,300italic,700,700italic" crossorigin="anonymous" referrerpolicy="no-referrer" />
	<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/normalize/8.0.1/normalize.min.css" integrity="sha512-NhSC1YmyruXifcj/KFRWoC561YpHpc5Jtzgvbuzx5VozKpWvQ+4nXhPdFgmx8xqexRcpAglTj9sIBWINXa8x5w==" crossorigin="anonymous" referrerpolicy="no-referrer" />
    <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/milligram/1.4.1/milligram.min.css" integrity="sha512-xiunq9hpKsIcz42zt0o2vCo34xV0j6Ny8hgEylN3XBglZDtTZ2nwnqF/Z/TTCc18sGdvCjbFInNd++6q3J0N6g==" crossorigin="anonymous" referrerpolicy="no-referrer" />
</head>
<body>
<section class="container">
    <table>
        <thead>
        <tr>
            <th>Project</th>
            <th>Description</th>
            <th>URL</th>
        </tr>
        </thead>
        <tbody>
        {{ range $relpath, $repo := . }}
        <tr>
            <td>{{ $repo.Name }}</td>
            <td>{{ $repo.Description }}</td>
            <td><a href="{{ $repo.URL }}">{{ $repo.URL }}</a></td>
        </tr>
        {{ end }}
        </tbody>
    </table>
</section>
</body>
</html>
`,
		),
)
