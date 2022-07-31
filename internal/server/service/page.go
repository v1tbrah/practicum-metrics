package service

import (
	"bytes"
	"text/template"
)

const pageTpl = `
<!DOCTYPE page>
<page>
	<head>
		<meta charset="UTF-8">
	</head>
	<body>
		<h1> Metrics </h1>
		<ul>
		{{range .Metrics}}<li>{{ . }}</li>{{else}}<div><strong>Have no information</strong></div>{{end}}
		</ul>
	</body>
</page>`

type dataForPage struct {
	Title   string
	Metrics []string
}

// NewDataForPage returns new dataForPage.
func NewDataForPage() *dataForPage {
	return &dataForPage{Title: "Data"}
}

// Page returns the completed start page template.
func (m *dataForPage) Page() (string, error) {
	buf := &bytes.Buffer{}
	t, err := template.New("webpage").Parse(pageTpl)
	if err != nil {
		return "", err
	}
	err = t.Execute(buf, m)
	if err != nil {
		return "", err
	}
	return buf.String(), nil

}
