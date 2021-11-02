# Kotone-DiVE

[![Go Report Card](https://goreportcard.com/badge/github.com/tpc3/kotone-dive)](https://goreportcard.com/report/github.com/tpc3/kotone-dive)
[![Docker Image CI](https://github.com/tpc3/Kotone-DiVE/actions/workflows/docker-image.yml/badge.svg)](https://github.com/tpc3/Kotone-DiVE/actions/workflows/docker-image.yml)
[![Go](https://github.com/tpc3/Kotone-DiVE/actions/workflows/go.yml/badge.svg)](https://github.com/tpc3/Kotone-DiVE/actions/workflows/go.yml)

In-development TTS bot for Discord, Re-written with golang.  
Suitable for self-host usage.

## Requirement

* Any computer (including raspberry pi or cloud hosting) that runs Windows / Mac / Linux
    * Linux is our main environment
* Discord bot account
* Token of your favorite TTS provider

## How-to

1. Install docker on your computer
    * ArchLinux: `pacman -Syu --needed docker`
1. [Download config](https://raw.githubusercontent.com/tpc3/Kotone-DiVE/master/config.yaml)
1. Edit the config file
    * Enter your bot token
    * Enable your TTS provider(s)
    * Disable TTS providers that you don't use
    * Adjust other settings
1. `docker run --rm -it -v $(PWD)/config.yaml:/Kotone-DiVE/config.yaml ghcr.io/tpc3/kotone-dive`
    * If you want to say "Docker sucks", just go ahead to the releases tab and download the binary
1. Profit

### Tip: database

If you want to save database on docker:

1. Change db path to `data/bbolt.db`
1. `docker run --rm -it -v $(PWD):/Kotone-DiVE/config.yaml -v $(PWD):/Kotone-DiVE/data ghcr.io/tpc3/kotone-dive`
1. Profit
