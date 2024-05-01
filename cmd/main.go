package main

import (
	"fmt"
	"github.com/gookit/slog"
	"github.com/takama/daemon"
	oneuptime_InfrastructureAgent_go "oneuptime-InfrastructureAgent-go"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"syscall"
)

const (
	appName = "oneuptime-infrastructure-agent"
	appDesc = "OneUptime Infrastructure Agent"
)

//var stdlog, errlog *log.Logger

type Service struct {
	daemon.Daemon
}

// Manage by daemon commands or run the daemon
func (service *Service) Manage() (string, error) {
	usage := "Usage: oneuptime-infrastructure-agent install | remove | start | stop | status"
	// If received any kind of command, do it
	if len(os.Args) > 1 {
		command := os.Args[1]
		switch command {
		case "install":
			return service.Install()
		case "remove":
			return service.Remove()
		case "start":
			return service.Start()
		case "stop":
			// No need to explicitly stop cron since job will be killed
			return service.Stop()
		case "status":
			return service.Status()
		default:
			return usage, nil
		}
	}
	// Set up channel on which to send signal notifications.
	// We must use a buffered channel or risk missing the signal
	// if we're not ready to receive when the signal is sent.
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, os.Kill, syscall.SIGTERM)
	sk := ""
	url := ""
	if len(os.Args) == 1 {
		return usage, nil
	}
	if os.Args[1] == "start" {
		if os.Args[2] == "" {
			return "secret-key is required flag, usage: --secret-key=<secret-key>", nil
		} else {
			//process secret key
			skKV := strings.Split(os.Args[2], "=")
			if len(skKV) == 2 {
				sk = skKV[1]
			} else {
				return "usage: --secret-key=abcabcabc", nil
			}
		}
		if os.Args[3] == "" {
			return "oneuptime-url is required flag, usage: --oneuptime-url=<url>", nil
		} else {
			//process oneuptime url
			urlKV := strings.Split(os.Args[3], "=")
			if len(urlKV) == 2 {
				url = urlKV[1]
			} else {
				return "usage: --oneuptime-url=https://oneuptime.example.org", nil
			}
		}
	}
	app := oneuptime_InfrastructureAgent_go.NewAgent(sk, url)
	app.Start()
	// Waiting for interrupt by system signal
	killSignal := <-interrupt
	slog.Info(fmt.Sprintf("Got signal: %s", killSignal))
	app.Close()
	return "Service exited", nil
}

func main() {
	//slog.Configure(func(l *slog.SugaredLogger) {
	//	l.Output = os.Stdout
	//})
	//fileHandler := handler.MustTimeRotateFile("oneuptime_agent.log", rotatefile.EveryDay)
	//slog.PushHandler(fileHandler)
	dType := daemon.SystemDaemon
	if runtime.GOOS == "darwin" {
		dType = daemon.UserAgent
	}
	srv, err := daemon.New(appName, appDesc, dType)
	if err != nil {
		slog.Error(err)
		os.Exit(1)
	}
	service := &Service{srv}
	status, err := service.Manage()
	service.GetTemplate()
	if err != nil {
		slog.Error(err)
		os.Exit(1)
	}
	srv.GetTemplate()
	fmt.Println(status)
}
