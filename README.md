# Kotone-DiVE

[![Go Report Card](https://goreportcard.com/badge/github.com/tpc3/kotone-dive)](https://goreportcard.com/report/github.com/tpc3/kotone-dive)
[![Docker Image CI](https://github.com/tpc3/Kotone-DiVE/actions/workflows/docker-image.yml/badge.svg)](https://github.com/tpc3/Kotone-DiVE/actions/workflows/docker-image.yml)
[![Go](https://github.com/tpc3/Kotone-DiVE/actions/workflows/go.yml/badge.svg)](https://github.com/tpc3/Kotone-DiVE/actions/workflows/go.yml)

TTS(Text-To-Speech) bot for Discord, Re-written with golang.  
Suitable for self-host usage.

## Requirements

* Any computer (including raspberry pi or cloud hosting) that runs Windows / Mac / Linux
  * Linux is our main environment
  * We don't check it, but other every platform golang supports should run it.
* Discord bot account
* Token of your favorite TTS provider
  * Currently supports: `azure`, `gcp`, `gtts`(for testing), `voicetext`, `voicevox`, `watson`, `coeiroink`, `aquestalk-proxy`
  * Every single voice type of the provider should work.

## How-to

### Configure

1. [Download config](https://raw.githubusercontent.com/tpc3/Kotone-DiVE/master/config-template.yaml) and rename to `config.yaml`
1. Edit the config file
    * Enter your bot token
    * Enable your TTS provider(s)
    * Disable TTS providers that you don't use
    * Adjust other settings

### Setup

* Install docker on your computer
  * Our docker image currently provides: `linux/386`, `linux/amd64`, `linux/arm64`, `linux/arm/v6`, `linux/arm/v7`
  * ArchLinux: `pacman -Syu --needed docker`
* Or if you just want to test or do not want to use docker, simply download binary from releases or GitHub Actions.
  * Our CI currently builds: `linux-386`, `linux-amd64`, `linux-arm`, `linux-arm64`, `darwin-amd64`, `darwin-arm64`, `windows-386`, `windows-amd64`

### Run

1. `docker run --rm -it -v $(PWD)/config.yaml:/Kotone-DiVE/config.yaml ghcr.io/tpc3/kotone-dive`
1. Profit

### Tip: database

If you want to save database on docker:

1. Change db path to `data/bbolt.db`
1. `docker run --rm -it -v $(PWD):/Kotone-DiVE/config.yaml -v $(PWD):/Kotone-DiVE/data ghcr.io/tpc3/kotone-dive`
1. Profit

### Contribution

Contributions are always welcome.  
(Please make issue or PR with English or Japanese)
