//package main
//
//import (
//	"fmt"
//	"github.com/gookit/slog"
//	"github.com/takama/daemon"
//	oneuptime_InfrastructureAgent_go "oneuptime-InfrastructureAgent-go"
//	"os"
//	"os/signal"
//	"runtime"
//	"strings"
//	"syscall"
//)
//
//const (
//	appName = "oneuptime-infrastructure-agent"
//	appDesc = "OneUptime Infrastructure Agent"
//)
//
////var stdlog, errlog *log.Logger
//
//type Service struct {
//	daemon.Daemon
//}
//
//// Manage by daemon commands or run the daemon
//func (service *Service) Manage() (string, error) {
//	usage := "Usage: oneuptime-infrastructure-agent install | remove | start | stop | status"
//	// If received any kind of command, do it
//	if len(os.Args) > 1 {
//		command := os.Args[1]
//		switch command {
//		case "install":
//			return service.Install()
//		case "remove":
//			return service.Remove()
//		case "start":
//			return service.Start()
//		case "stop":
//			// No need to explicitly stop cron since job will be killed
//			return service.Stop()
//		case "status":
//			return service.Status()
//		default:
//			return usage, nil
//		}
//	}
//	// Set up channel on which to send signal notifications.
//	// We must use a buffered channel or risk missing the signal
//	// if we're not ready to receive when the signal is sent.
//	interrupt := make(chan os.Signal, 1)
//	signal.Notify(interrupt, os.Interrupt, os.Kill, syscall.SIGTERM)
//	sk := ""
//	url := ""
//	if len(os.Args) == 1 {
//		return usage, nil
//	}
//	if os.Args[1] == "start" {
//		if os.Args[2] == "" {
//			return "secret-key is required flag, usage: --secret-key=<secret-key>", nil
//		} else {
//			//process secret key
//			skKV := strings.Split(os.Args[2], "=")
//			if len(skKV) == 2 {
//				sk = skKV[1]
//			} else {
//				return "usage: --secret-key=abcabcabc", nil
//			}
//		}
//		if os.Args[3] == "" {
//			return "oneuptime-url is required flag, usage: --oneuptime-url=<url>", nil
//		} else {
//			//process oneuptime url
//			urlKV := strings.Split(os.Args[3], "=")
//			if len(urlKV) == 2 {
//				url = urlKV[1]
//			} else {
//				return "usage: --oneuptime-url=https://oneuptime.example.org", nil
//			}
//		}
//	}
//	app := oneuptime_InfrastructureAgent_go.NewAgent(sk, url)
//	app.Start()
//	// Waiting for interrupt by system signal
//	killSignal := <-interrupt
//	slog.Info(fmt.Sprintf("Got signal: %s", killSignal))
//	app.Close()
//	return "Service exited", nil
//}
//
//func main() {
//	//slog.Configure(func(l *slog.SugaredLogger) {
//	//	l.Output = os.Stdout
//	//})
//	//fileHandler := handler.MustTimeRotateFile("oneuptime_agent.log", rotatefile.EveryDay)
//	//slog.PushHandler(fileHandler)
//	dType := daemon.SystemDaemon
//	if runtime.GOOS == "darwin" {
//		dType = daemon.UserAgent
//	}
//	srv, err := daemon.New(appName, appDesc, dType)
//	if err != nil {
//		slog.Error(err)
//		os.Exit(1)
//	}
//	service := &Service{srv}
//	status, err := service.Manage()
//	service.GetTemplate()
//	if err != nil {
//		slog.Error(err)
//		os.Exit(1)
//	}
//	srv.GetTemplate()
//	fmt.Println(status)
//}

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/gookit/config/v2"
	"github.com/gookit/slog"
	"github.com/kardianos/service"
	oneuptime_InfrastructureAgent_go "oneuptime-InfrastructureAgent-go"
	"oneuptime-InfrastructureAgent-go/logger"
	"os"
	"path/filepath"
	"runtime"
)

type configFile struct {
	SecretKey    string `json:"secret_key"`
	OneUptimeURL string `json:"oneuptime_url"`
}

func newConfigFile() *configFile {
	return &configFile{
		SecretKey:    "",
		OneUptimeURL: "",
	}
}

func (c *configFile) loadConfig() error {
	cfg := &configFile{}
	err := config.LoadFiles(c.configPath())
	if err != nil {
		return err
	}
	err = config.BindStruct("", cfg)
	if err != nil {
		return err
	}
	c.SecretKey = cfg.SecretKey
	c.OneUptimeURL = cfg.OneUptimeURL
	return nil
}

func (c *configFile) save(secretKey string, url string) error {
	err := c.loadConfig()
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	err = config.Set("secret_key", secretKey)
	if err != nil {
		return err
	}
	err = config.Set("oneuptime_url", url)
	if err != nil {
		return err
	}
	// Open the file with os.Create, which truncates the file if it already exists,
	// and creates it if it doesn't.
	file, err := os.Create(c.configPath())
	if err != nil {
		return err
	}
	defer file.Close()
	// Create a JSON encoder that writes to the file, and use Encode method
	// which will write the map to the file in JSON format.
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "    ") // Optional: makes the output more readable
	return encoder.Encode(config.Data())
}

// removeConfigFile deletes the configuration file.
func (c *configFile) removeConfigFile() error {

	// Check if the file exists before attempting to remove it.
	if _, err := os.Stat(c.configPath()); os.IsNotExist(err) {
		// File does not exist, return an error or handle it accordingly.
		return os.ErrNotExist
	}

	// Remove the file.
	err := os.Remove(c.configPath())
	if err != nil {
		// Handle potential errors in deleting the file.
		return err
	}

	return nil
}

// ensureDir checks if a directory exists and makes it if it does not.
func (c *configFile) ensureDir(dirName string) error {
	// Check if the directory exists
	info, err := os.Stat(dirName)
	if os.IsNotExist(err) {
		// Directory does not exist, create it
		return os.MkdirAll(dirName, 0755)
	}
	if err != nil {
		return err
	}
	if !info.IsDir() {
		// Exists but is not a directory
		return os.ErrExist
	}
	return nil
}

// configPath returns the full path to the configuration file,
// ensuring the directory exists or creating it if it does not.
func (c *configFile) configPath() string {
	var basePath string
	if runtime.GOOS == "windows" {
		basePath = os.Getenv("PROGRAMDATA")
		if basePath == "" {
			basePath = fmt.Sprintf("C:%sProgramData", string(filepath.Separator))
		}
	} else {
		basePath = fmt.Sprintf("%setc", string(filepath.Separator))
	}

	// Define the directory path where the configuration file will be stored.
	configDirectory := filepath.Join(basePath, "oneuptime_infrastructure_agent")

	// Ensure the directory exists.
	err := c.ensureDir(configDirectory)
	if err != nil {
		slog.Fatalf("Failed to create config directory: %v", err)
	}

	// Return the full path to the configuration file.
	return filepath.Join(configDirectory, "config.json")
}

type program struct {
	exit   chan struct{}
	agent  *oneuptime_InfrastructureAgent_go.Agent
	config *configFile
}

func (p *program) Start(s service.Service) error {
	// Start should not block. Do the actual work async.
	go p.run()
	return nil
}

func (p *program) run() {
	// Actual service code here.
	fmt.Println("Service is running")
}

func (p *program) Stop(s service.Service) error {
	// Clean up here
	return nil
}

func main() {
	config.WithOptions(config.WithTagName("json"))
	cfg := newConfigFile()

	svcConfig := &service.Config{
		Name:        "oneuptime-infrastructure-agent",
		DisplayName: "OneUptime Infrastructure Agent",
		Description: "The OneUptime Infrastructure Agent (Golang Version) is a lightweight, open-source agent that collects system metrics and sends them to the OneUptime platform. It is designed to be easy to install and use, and to be extensible.",
	}

	prg := &program{
		config: cfg,
	}

	s, err := service.New(prg, svcConfig)
	if err != nil {
		slog.Fatal(err)
		os.Exit(2)
	}

	// Set up the logger
	errs := make(chan error, 5)
	l, err := s.Logger(errs)
	if err != nil {
		slog.Fatal(err)
		os.Exit(2)
	}

	logHandler := logger.NewServiceSysLogHandler(l)
	slog.PushHandler(logHandler)

	if len(os.Args) > 1 {
		cmd := os.Args[1]
		switch cmd {
		case "install":
			installFlags := flag.NewFlagSet("install", flag.ExitOnError)
			secretKey := installFlags.String("secret-key", "", "Secret key (required)")
			oneuptimeURL := installFlags.String("oneuptime-url", "", "Oneuptime endpoint root URL (required)")
			err := installFlags.Parse(os.Args[2:])
			if err != nil {
				slog.Fatal(err)
				os.Exit(2)
			}
			prg.config.SecretKey = *secretKey
			prg.config.OneUptimeURL = *oneuptimeURL
			if prg.config.SecretKey == "" || prg.config.OneUptimeURL == "" {
				slog.Fatal("The --secret-key and --oneuptime-url flags are required for the 'install' command")
				os.Exit(2)
			}
			// save configuration
			err = prg.config.save(prg.config.SecretKey, prg.config.OneUptimeURL)
			if err != nil {
				slog.Fatal(err)
				os.Exit(2)
			}
			// Install the service
			if err := s.Install(); err != nil {
				slog.Fatal("Failed to install service: ", err)
				os.Exit(2)
			}
			fmt.Println("Service installed")
		case "start":
			err := prg.config.loadConfig()
			if os.IsNotExist(err) {
				slog.Fatal("Service configuration not found. Please install the service properly.")
				os.Exit(2)
			}
			if err != nil {
				slog.Fatal(err)
				os.Exit(2)
			}
			if err != nil || prg.config.SecretKey == "" || prg.config.OneUptimeURL == "" {
				slog.Fatal("Service configuration not found or is incomplete. Please install the service properly.")
				os.Exit(2)
			}
			err = s.Run()
			if err != nil {
				slog.Fatal(err)
				os.Exit(2)
			}
		case "uninstall", "stop", "restart":
			err := service.Control(s, cmd)
			if err != nil {
				slog.Fatal(err)
				os.Exit(2)
			}
			if cmd == "uninstall" {
				// remove configuration file
				err := prg.config.removeConfigFile()
				if err != nil {
					slog.Fatal(err)
					os.Exit(2)
				}
			}
		default:
			slog.Error("Invalid command")
			os.Exit(2)
		}
	} else {
		fmt.Println("Usage: oneuptime-infrastructure-agent install | uninstall | start | stop | restart")
	}
}
