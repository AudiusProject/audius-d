package statusbar

import (
	"github.com/AudiusProject/audius-d/pkg/conf"
	"github.com/caseymrm/menuet"
)

const (
	logo = "https://dl.dropboxusercontent.com/s/b6wothpryr0887o/Glyph_White%402x.png?dl=1"
)

func ContextToMenuState() *menuet.MenuState {
	ctxs, _ := conf.GetContexts()
	selectedCtxName, _ := conf.GetCurrentContextName()
	//selectedCtx, _ := conf.ReadOrCreateContextConfig()

	menuet.App().Children = func() []menuet.MenuItem {
		return ListContexts(ctxs, selectedCtxName)
	}

	return &menuet.MenuState{
		Image: logo,
	}
}

func ListContexts(ctxs []string, selected string) []menuet.MenuItem {
	var items []menuet.MenuItem
	for _, ctx := range ctxs {
		state := false
		if ctx == selected {
			state = true
		}
		item := menuet.MenuItem{
			Text:  ctx,
			State: state,
		}
		items = append(items, item)
	}
	return items
}
