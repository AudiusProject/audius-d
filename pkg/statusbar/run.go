package statusbar

import (
	"github.com/AudiusProject/audius-d/pkg/conf"
	"github.com/caseymrm/menuet"
)

func StatusBar() {
	menuet.App().SetMenuState(&menuet.MenuState{
		Image: logo,
	})

	menuet.App().Children = func() []menuet.MenuItem {
		return []menuet.MenuItem{
			{
				Text: "contexts",
				Children: func() []menuet.MenuItem {
					ctxs, _ := conf.GetContexts()
					selectedCtxName, _ := conf.GetCurrentContextName()
					return ListContexts(ctxs, selectedCtxName)
				},
			},
		}
	}
}

func RunStatusBar() {
	StatusBar()
	menuet.App().RunApplication()
}
