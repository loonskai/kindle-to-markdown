---
title: {{ .Title }}
authors: {{ .Authors }}
citation: {{ .Citation }}
---
# {{ .Title }}
{{ range .Sections }}
## {{ .Title }}
{{ range .Notes }}
{{ .Highlight }}
<i>{{ .Heading }}</i>
{{ if .Text }}<b>{{ .Text }}</b>{{ end }}
{{ end }}
{{ end }}
