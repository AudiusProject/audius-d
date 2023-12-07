package statusbar

import (
	"github.com/caseymrm/menuet"
)

func RunStatusBar() {
	menuet.App().SetMenuState(&menuet.MenuState{
		Image: AudiusLogo,
	})
	menuet.App().Label = "audius"

	defaultState()

	go asyncUpdateContexts()

	menuet.App().RunApplication()
}
