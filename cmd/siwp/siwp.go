package main

import (
	"os"

	"github.com/flaccid/slack-incoming-webhook-proxy/proxy"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

var (
	VERSION = "v0.0.0-dev"
)

func beforeApp(c *cli.Context) error {
	if c.GlobalBool("debug") {
		log.SetLevel(log.DebugLevel)
	}
	return nil
}

func main() {
	app := cli.NewApp()
	app.Name = "siwp"
	app.Version = VERSION
	app.Usage = "siwp"
	app.Action = start
	app.Before = beforeApp
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "webhook-url,t",
			Usage: "url of the slack webhook",
		},
		cli.BoolFlag{
			Name:  "debug,d",
			Usage: "run in debug mode",
		},
	}
	app.Run(os.Args)
}

func start(c *cli.Context) error {
	log.Debug("initialising")

	if len(c.String("webhook-url")) < 1 {
		log.Errorf("no webhook url provided")
		os.Exit(1)
	}

	proxy.Serve(c.String("webhook-url"), "8000")

	return nil
}
