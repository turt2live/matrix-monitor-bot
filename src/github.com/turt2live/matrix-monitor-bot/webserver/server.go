package webserver

import (
	"net/http"
	"github.com/turt2live/matrix-monitor-bot/config"
	"html/template"
	"github.com/sirupsen/logrus"
	"path"
	"github.com/turt2live/matrix-monitor-bot/matrix"
	"strings"
	"sort"
	"github.com/turt2live/matrix-monitor-bot/tracker"
	"fmt"
	"time"
	"github.com/turt2live/matrix-monitor-bot/util"
)

type ComparedDomain struct {
	Domain      string
	SendTime    string
	ReceiveTime string
	AverageTime string
	HasSend     bool
	HasReceive  bool
	Status      string // TODO: Replace with the HasSend/HasReceive and future IsOnline
	Description string
}

type CompareTemplateFields struct {
	SelfDomain      string
	Domains         []ComparedDomain
	FeaturedDomains []ComparedDomain
	RelativePath    string // Needed for the layout.html
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
		SelfDomain:      config.Get().Webserver.DefaultCompareDomain,
		Domains:         make([]ComparedDomain, 0),
		FeaturedDomains: make([]ComparedDomain, 0),
		RelativePath:    baseHref,
	}

	if fields.SelfDomain == "" {
		fields.SelfDomain = mxClient.Domain
	}

	us := tracker.GetDomain(mxClient.Domain)
	domainsToUse := make([]string, 0)
	domainsToUse = append(domainsToUse, config.Get().Webserver.DefaultCompareToDomains...)
	if len(domainsToUse) == 0 {
		domainsToUse = tracker.GetDomainsExcept(mxClient.Domain)

		for _, r := range us.GetRemotes() {
			exists := false
			for _, e := range domainsToUse {
				if e == r {
					exists = true
					break
				}
			}
			if !exists {
				domainsToUse = append(domainsToUse, r)
			}
		}

		sort.Strings(domainsToUse)
	}
	domainsToUse = append(domainsToUse, config.Get().Webserver.FeaturedCompareDomains...)

	handledDomains := make(map[string]bool)
	for _, domain := range domainsToUse {
		if handledDomains[domain] {
			continue
		}
		handledDomains[domain] = true

		remote := us.CompareTo(domain)
		avgTime := (time.Duration((remote.Send.Nanoseconds()+remote.Receive.Nanoseconds())/2) * time.Nanosecond).Truncate(time.Millisecond)
		description := fmt.Sprint(avgTime)
		status := "ok"
		if !remote.HasSend && !remote.HasReceive {
			status = "danger"
			description = "offline"
		} else if avgTime > config.WebWarnStatusThreshold || !remote.HasSend || !remote.HasReceive {
			status = "warn"
		}
		compDomain := ComparedDomain{
			Domain:      domain,
			SendTime:    fmt.Sprint(remote.Send.Truncate(time.Millisecond)),
			ReceiveTime: fmt.Sprint(remote.Receive.Truncate(time.Millisecond)),
			HasSend:     remote.HasSend,
			HasReceive:  remote.HasReceive,
			AverageTime: fmt.Sprint(avgTime),
			Status:      status,
			Description: description,
		}

		if util.StrArrayContains(config.Get().Webserver.FeaturedCompareDomains, domain) {
			fields.FeaturedDomains = append(fields.FeaturedDomains, compDomain)
		} else {
			fields.Domains = append(fields.Domains, compDomain)
		}
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
