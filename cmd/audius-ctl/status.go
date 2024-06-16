package main

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/AudiusProject/audius-d/pkg/conf"
	"github.com/AudiusProject/audius-d/pkg/health"
	"github.com/AudiusProject/audius-d/pkg/logger"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

type model struct {
	table table.Model
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			if m.table.Focused() {
				m.table.Blur()
			} else {
				m.table.Focus()
			}
		case "q", "ctrl+c":
			return m, tea.Quit
		case "enter":
			return m, tea.Batch(
				tea.Printf("Let's go to %s!", m.table.SelectedRow()[1]),
			)
		}
	}
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return baseStyle.Render(m.table.View()) + "\n"
}

type hcResult struct {
	Host          string
	HealthSummary health.NodeHealthSummary
	Error         error
}

type diskUsage struct {
	sizeBytes uint64
	usedBytes uint64
}

const (
	nodeCol int = iota
	typeCol
	upCol
	healthyCol
	chainCol
	websocketCol
	ipCol
	dbCol
	diskCol
	uptimeCol
	commentCol
)

const (
	diskUsageWarningThreshold float64 = 0.80
	diskUsageErrorThreshold   float64 = 0.90
	dbSizeWarningThreshold    uint64  = 2
)

var (
	noStatus     = "n/a"
	ignoreHealth bool
	statusCmd    = &cobra.Command{
		Use:          "status [host ...]",
		Short:        "Check health of configured nodes",
		SilenceUsage: true, // do not print --help text on failed node health
		RunE: func(cmd *cobra.Command, args []string) error {
			ctxConfig, err := conf.ReadOrCreateContextConfig()
			if err != nil {
				return err
			}

			var nodesToCheck map[string]conf.NodeConfig
			if len(args) == 0 {
				nodesToCheck = ctxConfig.Nodes
			} else {
				nodesToCheck, err = filterNodesFromContext(args, ctxConfig)
				if err != nil {
					return err
				}
			}

			var wg sync.WaitGroup
			resultsChan := make(chan hcResult, len(nodesToCheck))
			for host, config := range nodesToCheck {
				wg.Add(1)
				go func(h string, c conf.NodeConfig) {
					defer wg.Done()
					result, err := health.CheckNodeHealth(h, c)
					resultsChan <- hcResult{
						Host:          h,
						HealthSummary: result,
						Error:         err,
					}
				}(host, config)
			}

			go func() {
				wg.Wait()
				close(resultsChan)
			}()

			var results []hcResult
			for r := range resultsChan {
				results = append(results, r)
			}
			sort.Slice(results, func(i, j int) bool {
				return results[i].Host < results[j].Host
			})

			err = writeResultsToTable(results)
			if ctxConfig.Network.DeployOn == conf.Devnet {
				devnetHealth := health.CheckDevnetHealth()
				foundUnhealthy := false
				for _, h := range devnetHealth.Hosts {
					if !h.Healthy {
						foundUnhealthy = true
					}
					fmt.Printf("%s %t\n", h.Host, h.Healthy)
				}
				if err == nil && foundUnhealthy {
					err = logger.Error("Unhealthy devnet")
				}
			}
			return err
		},
	}
)

func init() {
	statusCmd.Flags().BoolVarP(&ignoreHealth, "ignore-health", "i", false, "Return non-zero only if nodes aren't up, ignoring health")
}

func writeResultsToTable(results []hcResult) error {
	columns := []table.Column{
		{Title: "Node", Width: 25},
		{Title: "Type", Width: 10},
		{Title: "Up", Width: 8},
		{Title: "Healthy", Width: 8},
		{Title: "Chain", Width: 10},
		{Title: "Websocket", Width: 10},
		{Title: "Client IP", Width: 10},
		{Title: "DB", Width: 10},
		{Title: "Disk", Width: 10},
		{Title: "Uptime", Width: 10},
		{Title: "Comment", Width: 10},
	}

	rows := []table.Row{}
	var unhealthyNode bool
	for _, res := range results {
		row := table.Row{
			res.Host,
			fmt.Sprint(res.HealthSummary.Type),
			fmt.Sprint(res.HealthSummary.Up),
			noStatus,
			noStatus,
			noStatus,
			noStatus,
			noStatus,
			noStatus,
			noStatus,
			fmt.Sprint(res.Error),
		}
		if !res.HealthSummary.Up {
			unhealthyNode = true
			rows = append(rows, row)
			continue
		}
		if !res.HealthSummary.Healthy && !ignoreHealth {
			unhealthyNode = true
		}

		row[healthyCol] = fmt.Sprint(res.HealthSummary.Healthy)
		if res.HealthSummary.Type == conf.Identity {
			rows = append(rows, row)
			continue
		}

		row[dbCol] = fmt.Sprint(res.HealthSummary.DatabaseSizeBytes)
		row[diskCol] = ""
		// diskUsage{
		// 	usedBytes: res.HealthSummary.DiskSpaceUsedBytes,
		// 	sizeBytes: res.HealthSummary.DiskSpaceSizeBytes,
		// }
		row[uptimeCol] = fmt.Sprint(time.Now().Sub(res.HealthSummary.BootTime))
		if res.HealthSummary.IPCheck {
			row[ipCol] = "matched"
		} else {
			row[ipCol] = "unmatched/error"
		}
		if res.Error == nil && len(res.HealthSummary.Errors) != 0 {
			row[commentCol] = strings.Join(res.HealthSummary.Errors, ",")
		}

		if res.HealthSummary.Type == conf.Discovery {
			var chainStatus string
			if res.HealthSummary.ChainHealthy {
				if !res.HealthSummary.ChainPortOpen {
					chainStatus = "Port 30300 unreachable"
				} else {
					chainStatus = "healthy"
				}
			} else {
				chainStatus = "unhealthy"
			}
			row[chainCol] = chainStatus

			wsStatus := "unreachable"
			if res.HealthSummary.WebsocketHealthy {
				wsStatus = "healthy"
			}
			row[websocketCol] = wsStatus

			rows = append(rows, row)
		} else {
			rows = append(rows, row)
			continue
		}
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(6),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(s)

	m := model{t}
	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
	// t := table.NewWriter()
	// t.SetStyle(table.StyleColoredMagentaWhiteOnBlack)
	// t.SetOutputMirror(os.Stdout)
	// t.AppendHeader(table.Row{
	// 	"Node",
	// 	"Type",
	// 	"Up",
	// 	"Healthy",
	// 	"Chain",
	// 	"Websocket",
	// 	"Client IP",
	// 	"DB",
	// 	"Disk",
	// 	"Uptime",
	// 	"Comment",
	// })

	// healthTransformer := text.Transformer(func(val interface{}) string {
	// 	switch fmt.Sprint(val) {
	// 	case "true", "healthy", "matched":
	// 		return text.FgGreen.Sprint(val)
	// 	case "n/a", "<nil>":
	// 		return text.FgHiBlack.Sprint(val)
	// 	default:
	// 		return text.FgRed.Sprint(val)
	// 	}
	// })
	// dbSizeTransformer := text.Transformer(func(val interface{}) string {
	// 	if fmt.Sprint(val) == "n/a" || fmt.Sprint(val) == "<nil>" {
	// 		return text.FgHiBlack.Sprint(val)
	// 	}

	// 	ival, ok := val.(uint64)
	// 	if !ok {
	// 		return text.FgRed.Sprint("NaN")
	// 	}
	// 	gb := ival / 1024 / 1024 / 1024
	// 	result := fmt.Sprintf("%d GB", gb)
	// 	if gb > dbSizeWarningThreshold {
	// 		return text.FgWhite.Sprint(result)
	// 	} else {
	// 		return text.FgRed.Sprint(result)
	// 	}
	// })
	// diskSizeTransformer := text.Transformer(func(val interface{}) string {
	// 	if fmt.Sprint(val) == "n/a" || fmt.Sprint(val) == "<nil>" {
	// 		return text.FgHiBlack.Sprint(val)
	// 	}

	// 	du, ok := val.(diskUsage)
	// 	if !ok {
	// 		return text.FgRed.Sprint("NaN")
	// 	}
	// 	ugb := du.usedBytes / 1024 / 1024 / 1024
	// 	sgb := du.sizeBytes / 1024 / 1024 / 1024
	// 	result := fmt.Sprintf("%d/%d GB", ugb, sgb)
	// 	ratio := float64(ugb) / float64(sgb)
	// 	if ratio < diskUsageWarningThreshold {
	// 		return text.FgWhite.Sprint(result)
	// 	} else if ratio < diskUsageErrorThreshold {
	// 		return text.FgYellow.Sprint(result)
	// 	} else {
	// 		return text.FgRed.Sprint(result)
	// 	}
	// })
	// uptimeTransformer := text.Transformer(func(val interface{}) string {
	// 	if fmt.Sprint(val) == "n/a" || fmt.Sprint(val) == "<nil>" {
	// 		return text.FgHiBlack.Sprint(val)
	// 	}

	// 	dur, ok := val.(time.Duration)
	// 	if !ok {
	// 		return text.FgRed.Sprint("NaN")
	// 	}
	// 	return text.FgWhite.Sprint(dur.Round(time.Second))
	// })

	// t.SetColumnConfigs([]table.ColumnConfig{
	// 	{
	// 		Name:        "Up",
	// 		Transformer: healthTransformer,
	// 	}, {
	// 		Name:        "Healthy",
	// 		Transformer: healthTransformer,
	// 	}, {
	// 		Name:        "Chain",
	// 		Transformer: healthTransformer,
	// 	}, {
	// 		Name:        "Websocket",
	// 		Transformer: healthTransformer,
	// 	}, {
	// 		Name:        "Client IP",
	// 		Transformer: healthTransformer,
	// 	}, {
	// 		Name:        "DB",
	// 		Transformer: dbSizeTransformer,
	// 	}, {
	// 		Name:        "Disk",
	// 		Transformer: diskSizeTransformer,
	// 	}, {
	// 		Name:        "Uptime",
	// 		Transformer: uptimeTransformer,
	// 	}, {
	// 		Name:        "Comment",
	// 		Transformer: healthTransformer,
	// 	},
	// })

	// var unhealthyNode bool
	// for _, res := range results {
	// 	row := table.Row{
	// 		res.Host,
	// 		res.HealthSummary.Type,
	// 		res.HealthSummary.Up,
	// 		noStatus,
	// 		noStatus,
	// 		noStatus,
	// 		noStatus,
	// 		noStatus,
	// 		noStatus,
	// 		noStatus,
	// 		res.Error,
	// 	}
	// 	if !res.HealthSummary.Up {
	// 		unhealthyNode = true
	// 		t.AppendRow(row)
	// 		continue
	// 	}
	// 	if !res.HealthSummary.Healthy && !ignoreHealth {
	// 		unhealthyNode = true
	// 	}

	// 	row[healthyCol] = res.HealthSummary.Healthy
	// 	if res.HealthSummary.Type == conf.Identity {
	// 		t.AppendRow(row)
	// 		continue
	// 	}

	// 	row[dbCol] = res.HealthSummary.DatabaseSizeBytes
	// 	row[diskCol] = diskUsage{
	// 		usedBytes: res.HealthSummary.DiskSpaceUsedBytes,
	// 		sizeBytes: res.HealthSummary.DiskSpaceSizeBytes,
	// 	}
	// 	row[uptimeCol] = time.Now().Sub(res.HealthSummary.BootTime)
	// 	if res.HealthSummary.IPCheck {
	// 		row[ipCol] = "matched"
	// 	} else {
	// 		row[ipCol] = "unmatched/error"
	// 	}
	// 	if res.Error == nil && len(res.HealthSummary.Errors) != 0 {
	// 		row[commentCol] = res.HealthSummary.Errors
	// 	}

	// 	if res.HealthSummary.Type == conf.Discovery {
	// 		var chainStatus string
	// 		if res.HealthSummary.ChainHealthy {
	// 			if !res.HealthSummary.ChainPortOpen {
	// 				chainStatus = "Port 30300 unreachable"
	// 			} else {
	// 				chainStatus = "healthy"
	// 			}
	// 		} else {
	// 			chainStatus = "unhealthy"
	// 		}
	// 		row[chainCol] = chainStatus

	// 		wsStatus := "unreachable"
	// 		if res.HealthSummary.WebsocketHealthy {
	// 			wsStatus = "healthy"
	// 		}
	// 		row[websocketCol] = wsStatus

	// 		t.AppendRow(row)
	// 	} else {
	// 		t.AppendRow(row)
	// 		continue
	// 	}
	// }

	// t.Render()
	if unhealthyNode {
		return fmt.Errorf("One or more health checks failed")
	}

	return nil
}
