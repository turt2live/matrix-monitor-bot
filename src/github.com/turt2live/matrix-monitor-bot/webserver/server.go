package webserver

import (
	"net/http"
	"github.com/turt2live/matrix-monitor-bot/config"
	"html/template"
	"github.com/sirupsen/logrus"
	"path"
	"github.com/turt2live/matrix-monitor-bot/matrix"
	"fmt"
	"github.com/turt2live/matrix-monitor-bot/metrics"
)

type ComparedDomain struct {
	Domain      string
	SendTime    string
	ReceiveTime string
}

type CompareTemplateFields struct {
	SelfDomain string
	Domains    []ComparedDomain
}

var mxClient *matrix.Client

func InitServer(mux *http.ServeMux, client *matrix.Client) {
	mxClient = client

	fs := http.FileServer(http.Dir(config.Runtime.WebContentDir))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	mux.HandleFunc("/", serveCompare)
}

func serveCompare(w http.ResponseWriter, r *http.Request) {
	layout := path.Join(config.Runtime.WebContentDir, "layout.html")
	file := path.Join(config.Runtime.WebContentDir, "compare.html")

	fields := CompareTemplateFields{
		SelfDomain: config.Get().Webserver.DefaultCompareDomain,
		Domains:    make([]ComparedDomain, 0),
	}

	if fields.SelfDomain == "" {
		fields.SelfDomain = mxClient.Domain
	}

	domainsToUse := config.Get().Webserver.DefaultCompareToDomains
	if len(domainsToUse) == 0 {
		domainsToUse = metrics.ListDomainsWithSendTimes(mxClient.Domain)
	}

	for _, domain := range domainsToUse {
		fields.Domains = append(fields.Domains, ComparedDomain{
			Domain:      domain,
			SendTime:    fmt.Sprint(metrics.CalculateSendTime(fields.SelfDomain, domain)),
			ReceiveTime: fmt.Sprint(metrics.CalculateSendTime(domain, fields.SelfDomain)),
		})
	}

	tmpl, err := template.ParseFiles(layout, file)
	if err != nil {
		logrus.Error(err)
		http.Error(w, "Failed to load template", http.StatusInternalServerError)
		return
	}

	err = tmpl.ExecuteTemplate(w, "layout", &fields)
	if err != nil {
		logrus.Error(err)
	}
}
