package app

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

type monitor struct {
	config Config
	client *http.Client
	seen   map[string]bool
}

type Feed struct {
	Data []struct {
		Title   string `xml:"title"`
		Summary string `xml:"summary"`
		Id      string `xml:"id"`
		Updated string `xml:"updated"`
		Link    struct {
			Href string `xml:"href,attr"`
		} `xml:"link"`
		Author struct {
			Name  string `xml:"name"`
			Email string `xml:"email"`
		} `xml:"author"`
	} `xml:"entry"`
}

var removeSummaryContentsRE = regexp.MustCompile(`(?mUis)<summary(.*)</summary>`)

func NewMonitor(c Config, client *http.Client) *monitor {
	return &monitor{
		config: c,
		client: client,
		seen:   map[string]bool{},
	}
}

func (m *monitor) Run(ctx context.Context) error {
	now := time.Now()          // time indicator of last run
	run := make(chan struct{}) // run now
	initialSort := sync.Once{} // run initial sorting

	go func() {
		run <- struct{}{}
		tick := time.NewTicker(time.Minute)
		for range tick.C {
			select {
			case run <- struct{}{}:
			case <-ctx.Done():
				fmt.Print("\r  ")
				return
			}
		}
	}()

	for {

		// ~~~~ Exit condition ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
		select {
		case <-ctx.Done():
			return nil
		case <-run:
		}

		// ~~~~ Working ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
		data, err := m.readFeed(context.Background())
		if err != nil {
			logrus.Errorf("Reading stackoverflow: %v\n", err)
			continue
		}

		data = removeSummaryContentsRE.ReplaceAll(data, nil)

		var items Feed
		if err := xml.Unmarshal(data, &items); err != nil {
			logrus.Errorf("XmlUnmarshaling: %v\n", err)
			continue
		}

		initialSort.Do(func() {
			sort.Slice(items.Data, func(i, j int) bool {
				return items.Data[i].Id < items.Data[j].Id
			})
		})

		fmt.Print(".")

		var id string
		var once sync.Once

		for _, v := range items.Data {
			id = strings.SplitAfter(v.Id, "q/")[1]

			if _, ok := m.seen[id]; ok {
				continue
			}

			once.Do(func() {
				fmt.Println("")
				fmt.Println(time.Since(now))
			})

			m.seen[id] = true

			logrus.Infof("%s [%s]\n", v.Id, v.Title)
			now = time.Now()

			m.open(v.Id)
		}
	}
}

func (m monitor) open(url string) {
	if !m.config.Open || len(m.seen) <= 30 {
		return
	}

	exec.Command("open", url).Run()
}

func (m monitor) readFeed(ctx context.Context) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, m.config.Feed(), nil)
	if err != nil {
		return nil, fmt.Errorf("get: %w", err)
	}

	resp, err := m.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("get: %w", err)
	}

	defer resp.Body.Close()
	buf, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("body: %w", err)
	}

	return buf, nil
}
