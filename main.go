package main

import (
	"flag"
	"io/ioutil"
	"os"
	"os/exec"
	"syscall"

	"github.com/sevlyar/go-daemon"
	"github.com/google/logger"
	"github.com/prashantgupta24/mac-sleep-notifier/notifier"
)

var killFlag = flag.Bool("kill", false, "kill TouchBarServer")
var daemonFlag = flag.Bool("bg", false, "run in the background")

var resetTouchBarOnWake = true

// Execute "pkill TouchBarServer" as root.
func killServer() {
	// Raise priviliges.
	if err := syscall.Setuid(0); err != nil {
		logger.Errorf("raising privs: setuid(0): %s", err)
		os.Exit(1)
	}
	// Kill the TouchBarServer process.
	cmd := exec.Command("/usr/bin/pkill", "-x", "TouchBarServer")
	out, err := cmd.CombinedOutput()
	if err != nil && len(out) > 0 {
		logger.Errorf("exec pkill TouchBarServer: %s", out);
	} else if err != nil {
		logger.Errorf("exec pkill TouchBarServer: %s", err);
	} else if len(out) > 0 {
		logger.Infof("exec pkill: %s", out)
	}
}

// Call ourselves, with a flag that means
// "not drop privs, then kill TouchBarServer as root"
func killServer2() {
	exe := os.Args[0]
	cmd := exec.Command(exe, "-kill")
	out, err := cmd.CombinedOutput()
	if err != nil && len(out) > 0 {
		logger.Errorf("exec pkill TouchBarServer: %s", out);
	} else if err != nil {
		logger.Errorf("exec pkill TouchBarServer: %s", err);
	} else if len(out) > 0 {
		logger.Infof("exec pkill: %s", out)
	}
}

func notificationListener() {
	logger.Infoln("starting sleep notifier")
	notifierCh := notifier.GetInstance().Start()

	for {
		select {
		case activity := <-notifierCh:
			if activity.Type == notifier.Awake {
				logger.Infoln("machine awake")
				if resetTouchBarOnWake {
					killServer2()
				}
			}
			if activity.Type == notifier.Sleep {
				logger.Infoln("machine sleeping")
			}
		}
	}
}

func main() {
	logger.Init("fix-touchbar", true, true, ioutil.Discard)

	// We might be called with the -kill flag.
	flag.Parse()
	if *killFlag {
		killServer()
		os.Exit(0)
	}

	// Drop setuid privs.
	uid := os.Getuid()
	if err := syscall.Setuid(uid); err != nil {
		logger.Errorf("dropping privs: setuid(%d): %s", uid, err);
		os.Exit(1)
	}

	// Daemonize?
	var cntxt *daemon.Context = nil
	if *daemonFlag {
		cntxt = &daemon.Context{}
		d, err := cntxt.Reborn()
		if err != nil {
			logger.Errorf("unable to run in the background: %s", err)
			return
		}
		if d != nil {
			// Parent.
			return
		}
	}
	defer func() {
		if cntxt != nil {
			cntxt.Release()
		}
	}()

	// Listen to sleep/wake notifications.
	go notificationListener()

	// Start tray.
	tray()
}
