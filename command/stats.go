package command

import (
	"bytes"
	"flag"
	"fmt"
	"regexp"
	"sort"
	"strings"
	"sync"
	"text/tabwriter"
	"time"

	"github.com/sloppyio/cli/pkg/api"
	"github.com/sloppyio/cli/ui"
)

// StatsCommand is a Command implementation that is used to display usage
// statistics about an entire project.
type StatsCommand struct {
	UI       ui.UI
	Projects api.ProjectsGetter
	Apps     api.AppsGetMetricer
	showAll  bool
}

// Help should return long-form help text.
func (c *StatsCommand) Help() string {
	helpText := `
Usage: sloppy stats [OPTIONS] PROJECT

  Displays usage statistics of running instances (memory, traffic)

Options:
  -a, --all     Show all instances (default shows just running instances)

Examples:

  sloppy stats letschat
`
	return strings.TrimSpace(helpText)
}

// Run should run the actual command with the given CLI instance and
// command-line args.
func (c *StatsCommand) Run(args []string) int {
	cmdFlags := newFlagSet("stats", flag.ContinueOnError)
	cmdFlags.BoolVar(&c.showAll, "a", false, "")
	cmdFlags.BoolVar(&c.showAll, "all", false, "")
	if err := cmdFlags.Parse(args); err != nil {
		c.UI.Error(err.Error())
		c.UI.Output("See 'sloppy stats --help'.")
		return 1
	}

	if cmdFlags.NArg() < 1 {
		return c.UI.ErrorNotEnoughArgs("stats", "", 1)
	}

	if strings.Contains(cmdFlags.Arg(0), "/") {
		c.UI.Error(fmt.Sprintf("invalid project path \"%s\". \n", cmdFlags.Arg(0)))
		return 1
	}

	project, _, err := c.Projects.Get(cmdFlags.Arg(0))
	if err != nil {
		c.UI.ErrorAPI(err)
		return 1
	}

	stats, err := c.collect(project)
	if err != nil {
		c.UI.ErrorAPI(err)
		return 1
	}

	if len(stats) == 0 {
		c.UI.Output("No apps running")
		return 1
	}

	var buf bytes.Buffer
	w := new(tabwriter.Writer)
	w.Init(&buf, 0, 8, 0, '\t', 0)
	fmt.Fprintf(w, "CONTAINER \t CPU %% \t MEM / LIMIT \t MEM %% \t NET I/O Extern \t NET I/O Intern\n")

	var keys []string
	var latest api.Timestamp
	for _, stat := range stats {
		if stat.Time.After(latest.Time) {
			latest = stat.Time
		}
	}

	for k := range stats {
		diff := latest.Time.Sub(stats[k].Time.Time)
		if (diff < 5*time.Second && diff > -5*time.Second) || c.showAll {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)

	for _, k := range keys {
		fmt.Fprintf(w, "%s\n", stats[k])
	}

	w.Flush()

	c.UI.Output(buf.String())

	return 0
}

// Synopsis should return a one-line, short synopsis of the command.
func (c *StatsCommand) Synopsis() string {
	return "Display metrics of a running app"
}

// Stat represents a container's stats.
type stat struct {
	App               string
	Service           string
	Time              api.Timestamp
	ID                string  // Container
	CPU               float64 // CPUPercentage
	Memory            float64
	MemoryLimit       float64
	InternalNetworkRx float64
	InternalNetworkTx float64
	ExternalNetworkRx float64
	ExternalNetworkTx float64
}

func (s *stat) String() string {
	return fmt.Sprintf("%s/%s-%s \t %.1f%% \t %s / %.f MiB \t %.1f%% \t %s / %s \t %s / %s",
		s.Service, s.App, s.ID[:6],
		s.CPU,
		humanByte(s.Memory), s.MemoryLimit,
		float64(s.Memory/(1<<20))/float64(s.MemoryLimit)*100,
		humanByte(s.ExternalNetworkRx), humanByte(s.ExternalNetworkTx),
		humanByte(s.InternalNetworkRx), humanByte(s.InternalNetworkTx),
	)
}

type Metric struct {
	app        *api.App
	metricName string
	service    string
	seriesName string
	time       time.Time
	value      float64
}

type Metrics []*Metric

type MetricFetchResult struct {
	app     *api.App
	service string
	metrics api.Metrics
	err     error
}

// Collect collects all statistics from all running apps, aggregate and merge them.
func (c *StatsCommand) collect(project *api.Project) (map[string]*stat, error) {
	quit := make(chan struct{})
	defer close(quit)

	var errors []error
	var result Metrics

	for r := range c.fetchAll(project, quit) {
		if r.err != nil {
			switch r.err.(type) {
			// This type of error was ignored in previous implementation.
			// However, this implementation does not handling multiple API
			// error responses in general.
			// TODO api errors should be handled downstream
			case *api.ErrorResponse:
				c.UI.ErrorAPI(r.err)
			default:
				errors = append(errors, r.err)
			}
			continue
		}
		result = append(result, c.toMetricSlice(r.app, r.service, r.metrics)...)
	}

	if len(errors) > 0 {
		return nil, errors[0]
	}
	return c.aggregate(result)
}

func (c *StatsCommand) fetchAll(project *api.Project, quit chan struct{}) <-chan MetricFetchResult {
	results := make(chan MetricFetchResult)
	wg := sync.WaitGroup{}
	for _, service := range project.Services {
		for _, app := range service.Apps {
			if !c.showAll && app.StatusCount("running") == 0 {
				continue
			}

			wg.Add(1)
			go func(app *api.App, name, service string) {
				metrics, _, err := c.Apps.GetMetrics(name, service, *app.ID)
				r := MetricFetchResult{app, service, metrics, err}
				select {
				case results <- r:
				case <-quit:
				}
				wg.Done()
			}(app, *project.Name, *service.ID)
		}
	}
	go func() {
		wg.Wait()
		close(results)
	}()
	return results
}

func (c *StatsCommand) aggregate(metrics []*Metric) (map[string]*stat, error) {
	stats := make(map[string]*stat)
	regex := regexp.MustCompile(`^(\w+)-([\w-]+).([a-z0-9-]{36}|([a-z/+]+)$)`)
	var p stat
	for _, metric := range metrics {
		matches := regex.FindStringSubmatch(metric.seriesName)
		if len(matches) == 0 || len(matches) < 3 {
			return nil, fmt.Errorf("invalid metric series name %q", metric.seriesName)
		}
		id := matches[3]

		s := &stat{
			App:               *metric.app.ID,
			Service:           metric.service,
			Time:              api.Timestamp{Time: metric.time},
			ID:                id,
			MemoryLimit:       float64(*metric.app.Memory),
			CPU:               p.CPU,
			Memory:            p.Memory,
			ExternalNetworkRx: p.ExternalNetworkRx,
			InternalNetworkRx: p.InternalNetworkRx,
			ExternalNetworkTx: p.ExternalNetworkTx,
			InternalNetworkTx: p.InternalNetworkTx,
		}

		switch metric.metricName {
		case "container_cpu_usage_percentage":
			s.CPU = metric.value
		case "container_memory_usage_bytes":
			s.Memory = metric.value
		case "container_network_receive_bytes_per_second":
			if strings.HasSuffix(metric.seriesName, "eth0") {
				s.ExternalNetworkRx = metric.value
				break
			}
			s.InternalNetworkRx = metric.value
		case "container_network_transmit_bytes_per_second":
			if strings.HasSuffix(metric.seriesName, "eth0") {
				s.ExternalNetworkTx = metric.value
				break
			}
			s.InternalNetworkTx = metric.value
		}
		stats[id] = s
		p = *s // merge next with previous one
	}
	return stats, nil
}

func (c *StatsCommand) toMetricSlice(app *api.App, serviceID string, metrics api.Metrics) Metrics {
	result := make(Metrics, 0)
	for metricName, series := range metrics {
		for seriesName, dataPoints := range series {
			for _, dataPoint := range dataPoints {
				result = append(result, &Metric{
					app:        app,
					metricName: metricName,
					service:    serviceID,
					seriesName: seriesName,
					time:       dataPoint.X.Time,
					value:      *dataPoint.Y,
				})
			}
		}
	}
	return result
}

// HumanByte returns a human-readable size.
func humanByte(size float64) string {
	var abbrs = []string{"B", "KiB", "MiB", "GiB", "TiB", "PiB"}

	i := 0
	for size >= 1024 {
		size = size / 1024
		i++
	}

	return fmt.Sprintf("%.3g %s", size, abbrs[i])
}
