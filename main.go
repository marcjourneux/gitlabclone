package main

import (
	logrus "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"os"
)

var (
	token      string
	group      string
	targetPath string
	level      string
	keyPath    string
	logr       *logrus.Logger
)

func main() {
	logr = logrus.New()
	app := &cli.App{
		Name:   "gitlabclone",
		Usage:  "clone all the gitlab projects and subprojects below a group",
		Action: clone,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "access-token",
				Aliases:     []string{"t"},
				Value:       "",
				Usage:       "gitlab access token",
				Destination: &token,
				Required:    true,
			},
			&cli.StringFlag{
				Name:        "group",
				Aliases:     []string{"g"},
				Value:       "",
				Usage:       "id of gitlab group",
				Destination: &group,
				Required:    true,
			},
			&cli.StringFlag{
				Name:        "ssh key relative",
				Aliases:     []string{"k"},
				Usage:       "relative user path for ssh key",
				Destination: &keyPath,
			},
			&cli.StringFlag{
				Name:        "destination path",
				Aliases:     []string{"d"},
				Usage:       "local path where to clone the project",
				Destination: &targetPath,
			},
			&cli.StringFlag{
				Name:        "log-level",
				Aliases:     []string{"l"},
				Value:       "Info",
				Usage:       "Log level (error/warning/info/debug/trace)",
				Destination: &level,
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		logr.Fatal(err)
	}
}

// Clone projects repo of the group
func clone(c *cli.Context) error {

	//logr = logrus.New()
	l, err := logrus.ParseLevel(level)
	if err != nil {
		return err
	}
	logr.SetLevel(l)

	VisitAndClone(token, group, keyPath, targetPath)
	logr.Info("Cloning OK")

	return nil
}
