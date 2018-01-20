# slack-incoming-webhook-tools

## Installation

```
git clone https://github.com/flaccid/slack-incoming-webhook-tools.git \
  "$GOPATH/src/github.com/flaccid/slack-incoming-webhook-tools"
cd "$GOPATH/src/github.com/flaccid/slack-incoming-webhook-tools"
go get ./...
```

## Usage

### siwp

A reverse proxy that lets your private clients send something to Slack without a token.

    $ siwp \
        --webhook-url https://hooks.slack.com/services/ABC123/XYZ456/789Abc \
        --listen-port 8900

### siwc

A CLI tool to easily post to a Slack incoming webhook.

Example:

    $ siwc \
        --payload '{"text":"This is a line of text.\nAnd this is another one."}' \
        --webhook-url https://hooks.slack.com/services/ABC123/XYZ456/789Abc

Using a template:

    $ siwc \
        -t '{"text": ":rancher: deploy of *{{.STACK_NAME}}* (<{{.BUILDKITE_BUILD_URL}}|{{.BUILDKITE_BUILD_NUMBER}}>) by {{.BUILDKITE_BUILD_CREATOR}} to *{{.STACK_ENV}}* succeeded.\n{{.BUILDKITE_MESSAGE}}"}' \
        -u https://hooks.slack.com/services/ABC123/XYZ456/789Abc

It is also possible to set environment variables instead of using cli options, see `siwc help`. Of course, you can set `--webhook-url` to that of your `siwp` reverse proxy instead.

### Docker

Run up `siwp` publishing the listen port locally:

    $ docker run -it \
        -e SLACK_WEBHOOK_URL="$SLACK_WEBHOOK_URL" \
        -p 8080:8080 \
          flaccid/slack-incoming-webhook-proxy:latest

Note: currently only `siwp` is built and published on Docker Hub.

## Building

TODO

## Upstream Documentation

- https://api.slack.com/incoming-webhooks

License and Authors
-------------------
- Author: Chris Fordham (<chris@fordham-nagy.id.au>)

```text
Copyright 2018, Chris Fordham

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
```
