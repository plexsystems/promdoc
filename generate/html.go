package generate

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
)

// HTML finds all rules at the given path and its
// subdirectories and generates an HTML table
// from the found rules.
func HTML(path string, input string) (string, error) {
	ruleGroups, err := getRuleGroups(path, input)
	if err != nil {
		return "", fmt.Errorf("get rule groups: %w", err)
	}

	const htmlTemplate = `<html>
<head>
  <title>Prometheus Alerts</title>
  <meta charset="utf-8" />
  <style>
  body{font-family:Microsoft Yahei,Helvetica,arial,sans-serif;font-size:14px;line-height:1.6;background-color:#fff;padding:30px;color:#516272;min-width:95%;margin:0 auto}body>*:first-child{margin-top:0!important}body>*:last-child{margin-bottom:0!important}a{color:#4183c4}a.absent{color:#c00}a.anchor{display:block;padding-left:30px;margin-left:-30px;cursor:pointer;position:absolute;top:0;left:0;bottom:0}h1,h2,h3,h4,h5,h6{margin:20px 0 10px;padding:0;font-weight:700;-webkit-font-smoothing:antialiased;cursor:text;position:relative}h1{font-size:28px;color:#2b3f52}h2{font-size:24px;border-bottom:1px solid #dde4e9;color:#2b3f52}h3{font-size:18px;color:#2b3f52}h4{font-size:16px;color:#2b3f52}h5{font-size:14px;color:#2b3f52}h6{color:#2b3f52;font-size:14px}p,blockquote,ul,ol,dl,li,table{margin:15px 0;color:#516272}body>h2:first-child{margin-top:0;padding-top:0}body>h1:first-child{margin-top:0;padding-top:0}body>h1:first-child+h2{margin-top:0;padding-top:0}body>h3:first-child,body>h4:first-child,body>h5:first-child,body>h6:first-child{margin-top:0;padding-top:0}a:first-child h1,a:first-child h2,a:first-child h3,a:first-child h4,a:first-child h5,a:first-child h6{margin-top:0;padding-top:0}h1 p,h2 p,h3 p,h4 p,h5 p,h6 p{margin-top:0}li p.first{display:inline-block}li{margin:0}ul,ol{padding-left:30px}ul :first-child,ol :first-child{margin-top:0}table{padding:0;border-collapse:collapse}table tr{border-top:1px solid #cccccc;background-color:#fff;margin:0;padding:0}table tr:nth-child(2n){background-color:#f8f8f8}table tr th{font-weight:700;border:1px solid #cccccc;margin:0;padding:6px 13px}table tr td{border:1px solid #cccccc;margin:0;padding:6px 13px}table tr th :first-child,table tr td :first-child{margin-top:0}table tr th :last-child,table tr td :last-child{margin-bottom:0}code{margin:0 2px;padding:0 5px;white-space:pre-wrap;word-break:break-all;display:block;border:1px solid #eaeaea;background-color:#f8f8f8;border-radius:3px}td code{margin:0;padding:0;white-space:pre}*{-webkit-print-color-adjust:exact}@media screen and (min-width: 914px){body{width:960px;margin:0 auto}}@media print{table,td{page-break-inside:avoid}td{word-wrap:break-word}}
  </style>
</head>
  
<body>
  <nav>
    <ul>
    {{- range . }}
      <li><a href="#{{.Name}}">{{.Name}}</a></li>
    {{- end }}
    </ul>
  </nav>
  {{range .}}
  <h2 id="{{.Name}}">{{.Name}}</h2>
  <table>
    <thead>
      <tr>
        <th>Name</th>
        <th>Summary</th>
        <th>Description</th>
        <th>Severity</th>
        <th>Expr</th>
        <th>For</th>
        <th>Runbook</th>
      </tr>
    </thead>
    <tbody>
      {{- range .Rules}}
      <tr>
        <td>{{.Alert}}</td>
        <td>{{.Annotations.summary}}</td>
        {{- if .Annotations.message}}
        <td>{{.Annotations.message}}</td>
        {{- else if .Annotations.description}}
        <td>{{.Annotations.description}}</td>
        {{- else}}
        <td></td>
        {{- end}}
        <td>{{.Labels.severity}}</td>
        <td><code>{{.Expr}}</code></td>
        <td>{{.For}}</td>
        <td><a href="{{.Annotations.runbook_url}}" target="blank">{{.Annotations.runbook_url}}</td>
      </tr>
    {{- end}}
    </tbody>
  </table>
  {{- end}}
</body>
</html>
`

	parsedTemplate, _ := template.New("template").Parse(htmlTemplate)
	var doc bytes.Buffer
	err = parsedTemplate.Execute(&doc, ruleGroups)
	if err != nil {
		log.Println("Error executing template :", err)
		return "", nil
	}

	return doc.String(), nil
}
