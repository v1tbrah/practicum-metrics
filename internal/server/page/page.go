package page

import (
	"bytes"
	"html/template"
)

type dataForHTML struct {
	Title   string
	Metrics []string
}

func NewData() *dataForHTML {
	return &dataForHTML{Title: "Metrics"}
}

const tpl = `
<!DOCTYPE page>
<page>
	<head>
		<meta charset="UTF-8">
	</head>
	<body>
		<h1> Metrics </h1>
		<ul>
		{{range .Metrics}}<li>{{ . }}</li>{{else}}<div><strong>Have no info</strong></div>{{end}}
		</ul>
	</body>
</page>`

func (m *dataForHTML) CompletedTpl() (string, error) {
	buf := &bytes.Buffer{}
	t, err := template.New("webpage").Parse(tpl)
	if err != nil {
		return "", err
	}
	err = t.Execute(buf, m)
	if err != nil {
		return "", err
	}
	return buf.String(), nil

}
