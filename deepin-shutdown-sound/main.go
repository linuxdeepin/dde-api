package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"pkg.deepin.io/dde/api/soundutils"
	"pkg.deepin.io/lib/log"
)

var logger = log.NewLogger("api/shutdown-sound")

func main() {
	handleSignal()

	canPlay, theme, event, err := soundutils.GetShutdownSound()
	if err != nil {
		logger.Warning("Get shutdown sound failed:", err)

		canPlay = true
		theme = "deepin"
		event = soundutils.EventShutdown
	}

	if !canPlay {
		return
	}

	err = doPlayShutdwonSound(theme, event)
	if err != nil {
		logger.Error("Play shutdown sound failed:", theme, event, err)
	}
}

func handleSignal() {
	var sigs = make(chan os.Signal, 2)
	signal.Notify(sigs, os.Kill, os.Interrupt)
	go func() {
		sig := <-sigs
		switch sig {
		case os.Kill, os.Interrupt:
			// Nothing to do
			logger.Debug("Recieve signal:", sig.String())
		}
	}()
}

func doPlayShutdwonSound(theme, event string) error {
	out, err := exec.Command("/usr/lib/deepin-api/sound-theme-player",
		theme, event).CombinedOutput()
	if err != nil {
		return fmt.Errorf("%v", string(out))
	}
	return nil
}
