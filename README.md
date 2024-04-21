# kindle-to-markdown
Transform your exported Kindle notes into markdown

## How to use
Save exported Kindle notebook to your computer as HTML file. The follow the steps: 
```
git clone https://github.com/loonskai/kindle-to-markdown
./convert \
  --source=file:///Users/johndoe/Downloads/Brave%20New%20World_%20(Original%20Classic%20Editions)%20-%20Notebook.html
  --output=BraveNewWorld.md
  --template=NoteTemplateExample.md
```

| Flag |  | Description |
|------|----------|-------------|
| `--source` | required | Full local path to an exported Kindle note HTML file |
| `--output` | required | Desired path for the output `.md` file |
| `--template` | required | Path to a template file that will be used for output `.md` |

## Output template
Template format uses standard Go template syntax.  

```go
type SourcePage struct {
  Title    string // Title of the book
  Authors  string // Authors of the book
  Citation string // Citation of the book
  Sections []Section // The list of the sections with highlights
}

type Section struct {
  Title string // Book section title
  Notes []NoteItem // A list of all notes under the section
}

type Note struct {
  Heading   string // Information about the highlight with its location
  Highlight string // Highlight
  Text      string // Your note attached to the highlight 
}
```

Example of `NoteTemplateExample.md`:
```md
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
```
