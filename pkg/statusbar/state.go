//go:build mac
// +build mac

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

			creators := []menuet.MenuItem{{
				Type: menuet.Separator,
			}}

			for name, creator := range ac.CreatorNodes {
				creatorCtx := CreatorContext{
					creator,
					name,
				}
				item := creatorCtx.MenuItem()
				creators = append(creators, item)
			}

			items = append(items, creators...)

			discovery := []menuet.MenuItem{{
				Type: menuet.Separator,
			}}

			for name, discoveryNode := range ac.DiscoveryNodes {
				discoveryCtx := DiscoveryContext{
					discoveryNode,
					name,
				}
				item := discoveryCtx.MenuItem()
				discovery = append(discovery, item)
			}

			items = append(items, discovery...)

			identity := []menuet.MenuItem{{
				Type: menuet.Separator,
			}}

			for name, identityService := range ac.IdentityService {
				identityCtx := IdentityContext{
					identityService,
					name,
				}
				item := identityCtx.MenuItem()
				identity = append(identity, item)
			}

			items = append(items, identity...)

			return items
		},
	}
	return item
}

type CreatorContext struct {
	conf.CreatorConfig
	containerName string
}

func (cc *CreatorContext) MenuItem() menuet.MenuItem {
	item := menuet.MenuItem{
		Text: cc.containerName,
		Clicked: func() {
			// restart node here
		},
	}
	return item
}

type DiscoveryContext struct {
	conf.DiscoveryConfig
	containerName string
}

func (cc *DiscoveryContext) MenuItem() menuet.MenuItem {
	item := menuet.MenuItem{
		Text: cc.containerName,
		Clicked: func() {
			// restart node here
		},
	}
	return item
}

type IdentityContext struct {
	conf.IdentityConfig
	containerName string
}

func (cc *IdentityContext) MenuItem() menuet.MenuItem {
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
			{Text: contexts},
		}
	}
}

func updateContexts(ctxs []string, selected string) {
	var items []menuet.MenuItem
	for _, ctx := range ctxs {
		fullCtx, _ := conf.ReadContext(ctx)
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
				Text:       "🟢 audius-d is running",
				FontWeight: menuet.WeightBold,
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
				Text:     contexts,
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
		time.Sleep(250 * time.Millisecond)
	}
}
