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
	domain     string
	apipath    string
	logr       *logrus.Logger
)

func main() {
	logr = logrus.New()
	apipath = "gitlab/v4"
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
				Name:        "api-root-domain",
				Aliases:     []string{"r"},
				Usage:       "gitlab domain (default gitlab.com)", //https://dev.truckonline.pro/gitlab/api/v4
				Value:       "gitlab.com",
				Destination: &domain,
			},
			&cli.StringFlag{
				Name:        "api-path",
				Aliases:     []string{"a"},
				Usage:       "path to api for gitlab", //https://dev.truckonline.pro/gitlab/api/v4
				Value:       "api/v4",
				DefaultText: "api/v4",
				Destination: &apipath,
			},
			&cli.StringFlag{
				Name:        "group",
				Aliases:     []string{"g"},
				Value:       "",
				Usage:       "id of gitlab group",
				Destination: &group,
			},
			&cli.StringFlag{
				Name:        "sshkey-relative-path",
				Aliases:     []string{"k"},
				Usage:       "relative user path for ssh key",
				Destination: &keyPath,
			},
			&cli.StringFlag{
				Name:        "destination-path",
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
	logr.Info("domain :" + domain)
	VisitAndClone(token, domain, apipath, group, keyPath, targetPath)
	logr.Info("Cloning OK")

	return nil
}
