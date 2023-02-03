package main

import (
	"errors"
	gogit "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport"
	gitssh "github.com/go-git/go-git/v5/plumbing/transport/ssh"
	logrus "github.com/sirupsen/logrus"
	gitlab "github.com/xanzy/go-gitlab"
	"net/url"
	"os"
	"os/user"
	"path"
	"strconv"
	"strings"
	"sync"
)

// return the auth from the ssh private key
func auth(privateKey string, password string) (transport.AuthMethod, error) {
	return gitssh.NewPublicKeysFromFile("git", privateKey, password)
}

// Check if a file exists
func exists(path string) bool {
	_, err := os.Stat(path)
	return !errors.Is(err, os.ErrNotExist)
}

// visit all subgroup inside group
func VisitGroup(auth transport.AuthMethod, glClient *(gitlab.Client), groupID, targetPath string) {
	var groups []*gitlab.Group
	var wg sync.WaitGroup
	var err error
	logr.Infof("analysing group " + groupID + "\n")
	if len(groupID) == 0 {
		opt := &gitlab.ListGroupsOptions{}
		groups, _, err = glClient.Groups.ListGroups(opt)

		if err != nil {
			logr.Infof("error when listing groups err.\n")
			logr.Fatal(err)
		}
	} else {
		opt := &gitlab.ListSubgroupsOptions{}
		groups, _, err = glClient.Groups.ListSubgroups(groupID, opt)

		if err != nil {
			logr.Infof("error when listing subgroups err.\n")
			logr.Fatal(err)
		}
		logr.Infof("nb subgroup " + strconv.Itoa(len(groups)) + "\n")
	}

	for _, grp := range groups {
		// Create the needed group folder
		groupTargetPath := path.Join(targetPath, grp.FullPath)
		logr.Infof("visit group " + grp.FullPath + "\n")
		os.MkdirAll(groupTargetPath, os.ModePerm)

		VisitGroup(auth, glClient, strconv.Itoa(grp.ID), targetPath)

		wg.Add(1)
		go CloneProjects(auth, glClient, strconv.Itoa(grp.ID), groupTargetPath, &wg)
	}

	wg.Wait()

}

// Clone or update a project
func CloneUpdateProject(prj *(gitlab.Project), auth transport.AuthMethod, glClient *(gitlab.Client), groupTargetPath string) {
	logr.WithFields(logrus.Fields{
		"repository": prj.PathWithNamespace,
		"action":     "validating",
	})

	repoPath := path.Join(groupTargetPath, prj.Path)

	// test if the git repo is already there
	if exists(path.Join(repoPath, ".git")) {
		// does it have the correct remote
		repo, err := gogit.PlainOpen(repoPath)
		if err != nil {
			logr.WithFields(logrus.Fields{
				"repository": repoPath,
				"action":     "open",
				"error":      err,
			}).Errorf(" git repo PlainOpen failed\n")
			return
		}
		remotes, err := repo.Remotes()
		if err != nil {
			logr.WithFields(logrus.Fields{
				"repository": repoPath,
				"action":     "git -v remote",
				"error":      err,
			}).Errorf("Cannot get remote on repo \n")
			return
		}
		for _, remote := range remotes {
			// if exists get the fetch url and update if valid
			if len(remote.Config().URLs) >= 1 {
				repoURL := remote.Config().URLs[0]
				if s := strings.Split(repoURL, ":"); len(s) >= 2 {
					if strings.Contains(s[1], prj.PathWithNamespace) {
						wtree, err := repo.Worktree()
						if err != nil {
							logr.WithFields(logrus.Fields{
								"repository": repoPath,
								"action":     "wortree",
								"error":      err,
							}).Errorf("Cannot get repo worktree\n")
							return
						}
						err = wtree.Pull(&gogit.PullOptions{Auth: auth})
						if err != nil {
							logr.WithFields(logrus.Fields{
								"repository": repoPath,
								"action":     "update",
								"error":      err,
							}).Infof("Cannot update repository \n")
						} else {
							logr.WithFields(logrus.Fields{
								"repository": repoPath,
								"action":     "update",
								"error":      "OK",
							}).Infof("Repo correctly updated \n")
						}
					}
				}
			}
		}
		// |Not already a git folder, let's clone the repo
	} else {
		option := &gogit.CloneOptions{Auth: auth, URL: prj.SSHURLToRepo, Progress: logr.Out}
		_, err := gogit.PlainClone(repoPath, false, option)
		if err != nil {
			if !errors.Is(err, gogit.ErrRepositoryAlreadyExists) {
				logr.WithFields(logrus.Fields{
					"repository": repoPath,
					"action":     "clone",
					"error":      err,
				}).Error()
			}
		}
	}
}

// clone all projects in a given group
func CloneProjects(auth transport.AuthMethod, glClient *(gitlab.Client), groupID, groupTargetPath string, wg *sync.WaitGroup) {

	defer wg.Done()

	option := &gitlab.ListGroupProjectsOptions{
		ListOptions: gitlab.ListOptions{
			PerPage: 10,
			Page:    1,
		}}

	// Get paginated results
	for {
		// Get the first page with projects.
		projects, resp, err := glClient.Groups.ListGroupProjects(groupID, option)
		if err != nil {
			logr.WithFields(logrus.Fields{
				"group":  groupID,
				"action": "list project",
				"error":  err,
			}).Errorf("Cannot list group projects from gitlab")
			return
		}

		// List all the projects we've found so far.
		for _, prj := range projects {
			if !prj.Archived {
				CloneUpdateProject(prj, auth, glClient, groupTargetPath)
			}
		}

		// Exit the loop when we've seen all pages.
		if resp.NextPage == 0 {
			break
		}

		// Update the page number to get the next page.
		option.Page = resp.NextPage
	}
}

// GroupCloneAllProjects clones all gitlab projects in the given group and or project and it's subgroups
func VisitAndClone(token, domain, apipath, groupID, keyPath, targetPath string) {

	// Get ssh key
	usr, err := user.Current()
	if err != nil {
		logr.Fatal(err)
	}
	auth, err := auth(path.Join(usr.HomeDir, keyPath), "")
	if err != nil {
		logr.Fatal(err)
	}
	// Create the gitlab client
	apiurl := &url.URL{
		Scheme: "https",
		Host:   domain,
		Path:   apipath,
	}
	logr.WithFields(logrus.Fields{
		"action": "create gitlab client",
	}).Infof(apiurl.String())
	glClient, err := gitlab.NewClient(token, gitlab.WithBaseURL(apiurl.String()))
	if err != nil {
		logr.WithFields(logrus.Fields{
			"action": "get gitlab client",
			"error":  err,
		}).Fatalf("Failed to create gitlab client")
	}

	VisitGroup(auth, glClient, groupID, targetPath)

}
