package main

import (
	"github.com/codegangsta/cli"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Name = "httpdiy"
	app.Version = "v0.0.1"
	app.Authors = []cli.Author{
		cli.Author{
			Name:  "argent_smith",
			Email: "argentoff@evrone.ru",
		},
	}
	app.Usage = "HTTP-do-it-yourself testing server for a Vexor cache tool"

	app.Commands = []cli.Command{
		cli.Command{
			Name:        "server",
			Aliases:     []string{"s"},
			Usage:       "start the server",
			Description: "The tiny web-server runs in foreground on port 4242",
			Action:      runServer,
		},
	}

	app.Run(os.Args)
}
