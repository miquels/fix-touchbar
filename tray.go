package main

import (
	_ "embed"

	"github.com/getlantern/systray"
	"github.com/google/logger"
)

//go:embed resources/Pictogrammers-Material-Table-row.512.png
var iconData []byte

func tray() {
	onExit := func() {
		logger.Infoln("quit")
	}
	systray.Run(onReady, onExit)
}

func onReady() {
	// Enable tray
	systray.SetTemplateIcon(iconData, iconData)
	systray.SetTooltip("Fix Touchbar")

	// Build the menu.
	mChecked := systray.AddMenuItemCheckbox("Reset touchbar after every wake-up", "reset on wake", resetTouchBarOnWake)
	mReset := systray.AddMenuItem("Reset Touchbar NOW", "reset now")
	systray.AddSeparator()
	mQuit := systray.AddMenuItem("Quit", "Quit fix-touchbar")

	for {
		select {
		case <-mChecked.ClickedCh:
			if mChecked.Checked() {
				mChecked.Uncheck()
			} else {
				mChecked.Check()
			}
			resetTouchBarOnWake = mChecked.Checked()
		case <-mReset.ClickedCh:
			killServer2()
		case <-mQuit.ClickedCh:
			systray.Quit()
			break
		}
	}
}
