package main

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig"
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

// returns key, value for all environment variables starting with prefix
func environment(prefix string) map[string]string {
	env := make(map[string]string)
	for _, setting := range os.Environ() {
		pair := strings.SplitN(setting, "=", 2)
		if strings.HasPrefix(pair[0], prefix) {
			env[pair[0]] = pair[1]
		}
	}
	return env
}

func customFuncMap() template.FuncMap {
	var functionMap = map[string]interface{}{"environment": environment}
	return template.FuncMap(functionMap)
}

func parse(s string) (*template.Template, error) {
	p, err := template.New("").Funcs(sprig.TxtFuncMap()).Funcs(customFuncMap()).Parse(s)
	log.Debug("P", p)

	return p, err
}

func readEnv() (env map[string]string) {
	env = make(map[string]string)
	for _, setting := range os.Environ() {
		pair := strings.SplitN(setting, "=", 2)
		env[pair[0]] = pair[1]
	}
	return
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
		env := readEnv()
		log.Debug("env: ", env)

		// load template
		t, err := parse(c.String("template"))
		if err != nil {
			panic(err)
		}

		// render the template
		var tpl bytes.Buffer
		if err := t.Execute(&tpl, env); err != nil {
			panic(err)
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
		panic(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	log.WithFields(log.Fields{
		"status":  resp.Status,
		"headers": resp.Header,
		"body":    string(body),
	}).Debug("response")

	return nil
}
