package main

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/flaccid/slack-incoming-webhook-tools/util"
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
	app.Name = "siwc"
	app.Version = VERSION
	app.Usage = "wrapper to easily post to a slack incoming webhook"
	app.Action = start
	app.Before = beforeApp
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "webhook-url,u",
			Usage:  "url of the slack webhook",
			EnvVar: "SLACK_WEBHOOK_URL",
		},
		cli.StringFlag{
			Name:   "payload,p",
			Usage:  "json payload",
			EnvVar: "SLACK_WEBHOOK_PAYLOAD",
		},
		cli.StringFlag{
			Name:   "template,t",
			Usage:  "payload template",
			EnvVar: "SLACK_WEBHOOK_PAYLOAD_TEMPLATE",
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

	if len(c.String("payload")) < 1 && len(c.String("template")) < 1 {
		log.Errorf("no webhook payload or template provided")
		os.Exit(1)
	}

	url := c.String("webhook-url")
	var payload []byte

	// if a template is provided, render it to payload
	if len(c.String("template")) < 1 {
		payload = []byte(c.String("payload"))
	} else {
		// get environment variables to supply to the template
		env := util.ReadEnv()
		log.Debug("env: ", env)

		// load template
		t, err := util.Parse(c.String("template"))
		if err != nil {
			log.Fatal(err)
		}

		// render the template
		var tpl bytes.Buffer
		if err := t.Execute(&tpl, env); err != nil {
			log.Fatal(err)
		}
		log.Debug("rendered: ", tpl.String())

		payload = []byte(tpl.String())
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))

	req.Header.Set("User-Agent", "siwc/"+VERSION)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	log.WithFields(log.Fields{
		"payload": string(payload),
		"headers": req.Header,
	}).Debug("request")

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	log.WithFields(log.Fields{
		"status":  resp.Status,
		"headers": resp.Header,
		"body":    string(body),
	}).Debug("response")

	return nil
}
