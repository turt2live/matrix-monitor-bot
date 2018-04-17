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
	"time"
	"strings"
	"sort"
)

type ComparedDomain struct {
	Domain      string
	SendTime    string
	ReceiveTime string
	AverageTime string
	Status      string
	Description string
}

type CompareTemplateFields struct {
	SelfDomain   string
	Domains      []ComparedDomain
	RelativePath string // Needed for the layout.html
}

var mxClient *matrix.Client
var baseHref string

func InitServer(mux *http.ServeMux, client *matrix.Client) {
	mxClient = client

	fs := http.FileServer(http.Dir(config.Runtime.WebContentDir))
	prefix := config.Get().Webserver.RelativePath
	if !strings.HasSuffix(prefix, "/") {
		prefix = prefix + "/"
	}
	if !strings.HasPrefix(prefix, "/") {
		prefix = "/" + prefix
	}
	baseHref = prefix
	prefix = prefix + "static/"
	mux.Handle(prefix, http.StripPrefix(prefix, fs))

	mux.HandleFunc("/", serveCompare)
}

func serveCompare(w http.ResponseWriter, r *http.Request) {
	layout := path.Join(config.Runtime.WebContentDir, "layout.html")
	file := path.Join(config.Runtime.WebContentDir, "compare.html")

	fields := CompareTemplateFields{
		SelfDomain:   config.Get().Webserver.DefaultCompareDomain,
		Domains:      make([]ComparedDomain, 0),
		RelativePath: baseHref,
	}

	if fields.SelfDomain == "" {
		fields.SelfDomain = mxClient.Domain
	}

	domainsToUse := config.Get().Webserver.DefaultCompareToDomains
	if len(domainsToUse) == 0 {
		domainsToUse = metrics.ListDomainsWithSendTimes(mxClient.Domain)
		sort.Strings(domainsToUse)
	}

	for _, domain := range domainsToUse {
		sendTime := metrics.CalculateSendTime(fields.SelfDomain, domain).Truncate(time.Millisecond)
		receiveTime := metrics.CalculateSendTime(domain, fields.SelfDomain).Truncate(time.Millisecond)
		avgTime := (time.Duration((sendTime.Nanoseconds()+receiveTime.Nanoseconds())/2.0) * time.Nanosecond).Truncate(time.Millisecond)
		description := fmt.Sprint(avgTime)
		status := "ok"
		if avgTime == 0 {
			status = "danger"
			description = "offline"
		} else if avgTime > config.WebWarnStatusThreshold {
			status = "warn"
		}
		fields.Domains = append(fields.Domains, ComparedDomain{
			Domain:      domain,
			SendTime:    fmt.Sprint(sendTime),
			ReceiveTime: fmt.Sprint(receiveTime),
			AverageTime: fmt.Sprint(avgTime),
			Status:      status,
			Description: description,
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
