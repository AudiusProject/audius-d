//go:build osx
// +build osx

package statusbar

import (
	"os/exec"
	"time"

	"github.com/AudiusProject/audius-d/pkg/conf"
	"github.com/caseymrm/menuet"
)

type AudiusContext struct {
	conf.ContextConfig
	name string
}

func (ac *AudiusContext) MenuItem(selectedctx string) menuet.MenuItem {
	item := menuet.MenuItem{
		Text:  ac.name,
		State: ac.name == selectedctx,
		Clicked: func() {
			if ac.name == selectedctx {
				return
			}
			conf.UseContext(ac.name)
		},
		Children: func() []menuet.MenuItem {
			items := []menuet.MenuItem{
				{
					Text:       "configured nodes",
					FontWeight: menuet.WeightRegular,
				},
				{
					Text:       "edit configuration",
					FontWeight: menuet.WeightRegular,
					Clicked: func() {
						confDir, _ := conf.GetContextBaseDir()
						exec.Command("code", confDir).Run()
					},
				},
			}

			nodes := []menuet.MenuItem{{
				Type: menuet.Separator,
			}}

			for name, node := range ac.ContextConfig.Nodes {
				creatorCtx := NodeContext{
					node,
					name,
				}
				item := creatorCtx.MenuItem()
				nodes = append(nodes, item)
			}

			items = append(items, nodes...)

			return items
		},
	}
	return item
}

type NodeContext struct {
	conf.NodeConfig
	containerName string
}

func (cc *NodeContext) MenuItem() menuet.MenuItem {
	item := menuet.MenuItem{
		Text: cc.containerName,
		Clicked: func() {
			// restart node here
		},
	}
	return item
}

func defaultState() {
	menuet.App().Children = func() []menuet.MenuItem {
		return []menuet.MenuItem{
			{Text: "contexts"},
		}
	}
}

func updateContexts(ctxs []string, selected string) {
	var items []menuet.MenuItem
	for _, ctx := range ctxs {
		fullCtx, _ := conf.ReadOrCreateContextConfig()
		ac := AudiusContext{
			*fullCtx,
			ctx,
		}
		item := ac.MenuItem(selected)
		items = append(items, item)
	}

	menuet.App().Children = func() []menuet.MenuItem {
		return []menuet.MenuItem{
			{
				Text: "🟢 audius-d is running",
			},
			{
				Type: menuet.Separator,
			},
			{
				Text:       "Dashboard",
				FontWeight: menuet.WeightRegular,
				Clicked: func() {
					url := "https://dashboard.audius.org"
					exec.Command("open", url).Start()
				},
			},
			{
				Text:       "Healthz",
				FontWeight: menuet.WeightRegular,
				Clicked: func() {
					url := "https://healthz.audius.co"
					exec.Command("open", url).Start()
				},
			},
			{
				Text:       "Listen",
				FontWeight: menuet.WeightRegular,
				Clicked: func() {
					url := "https://audius.co"
					exec.Command("open", url).Start()
				},
			},
			{
				Type: menuet.Separator,
			},
			{
				Text:     "contexts",
				Children: func() []menuet.MenuItem { return items },
			},
		}
	}
	menuet.App().MenuChanged()
}

// polls the context info and updates the chart accordingly
func asyncUpdateContexts() {
	for {
		ctxs, _ := conf.GetContexts()
		selectedctx, _ := conf.GetCurrentContextName()

		updateContexts(ctxs, selectedctx)
		time.Sleep(2000 * time.Millisecond)
	}
}
