package service

import (
	"bytes"
	"text/template"

	"github.com/rs/zerolog/log"
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
		{{range .Data}}<li>{{ . }}</li>{{else}}<div><strong>Have no information</strong></div>{{end}}
		</ul>
	</body>
</page>`

type dataForPage struct {
	Title string
	Data  []string
}

// NewDataForPage returns new dataForPage.
func NewDataForPage() *dataForPage {
	log.Debug().Msg("service.NewDataForPage started")
	defer log.Debug().Msg("service.NewDataForPage ended")

	return &dataForPage{Title: "Metrics"}
}

// Page returns the completed start page template.
func (m *dataForPage) Page() (string, error) {
	log.Debug().Msg("service.Page started")
	defer log.Debug().Msg("service.Page ended")

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
