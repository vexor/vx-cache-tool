package main

import (
	"github.com/codegangsta/cli"
	"os"
	"path/filepath"
)

// Sets up application's configuration
func prepareApp(c *cli.Context) (err error) {
	cfg.build(c)
	err = os.MkdirAll(cfg.BackupDir, 0755)
	check(err)
	mtimes.restore()
	return
}

func main() {
	app := cli.NewApp()
	app.Name = "vx-cache-tool"
	app.Version = version()
	app.Authors = authors()
	app.Usage = "Vexor(tm) cache management tool"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "config, f",
			Value: "vx-cache-tool.cfg",
			Usage: "config file",
		},

		cli.StringFlag{
			Name:   "backup_dir, d",
			Usage:  "directory for the tool's own technical backups",
			Value:  filepath.Join(os.Getenv("HOME"), ".vx-cache-tool"),
			EnvVar: "VX_CACHE_TOOL_DIR",
		},
	}

	app.Before = prepareApp

	app.Commands = []cli.Command{
		{
			Name:    "add",
			Aliases: []string{"a"},
			Usage:   "adds listed paths to cache",
			Action:  doAdd,
		},
		{
			Name:    "fetch",
			Aliases: []string{"f"},
			Usage:   "fetching an archive from a specified url",
			Action:  doFetch,
		},
		{
			Name:    "push",
			Aliases: []string{"p"},
			Usage:   "pushes cache archive to the specified url",
			Action:  doPush,
		},
	}

	app.Run(os.Args)
}
