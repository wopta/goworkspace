package common

import (
	"fmt"
	"os"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
)

func CloneRepo() (*git.Repository, string) {
	fmt.Printf("Cloning of repo %s...\n", os.Getenv("GIT_REMOTE"))
	dir, err := os.MkdirTemp("", "gowork")
	if err != nil {
		panic(err)
	}
	repo, err := git.PlainClone(dir, false, &git.CloneOptions{
		URL: os.Getenv("GIT_REMOTE"),
		Auth: &http.BasicAuth{
			Username: "ci-bot",
			Password: os.Getenv("GIT_ACCESS_TOKEN"),
		},
		ReferenceName: "master",
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("Clone completed!")
	return repo, dir
}

func CreateTag(repo *git.Repository, newTag, message string, hash plumbing.Hash) {
	fmt.Printf("Creating tag %s...\n", newTag)
	if _, err := repo.CreateTag(newTag, hash, &git.CreateTagOptions{
		Message: message,
		Tagger: &object.Signature{
			Name:  "Cloudbuild CI",
			Email: "technology@wopta.it",
			When:  time.Now(),
		},
	}); err != nil {
		panic(err)
	}
	fmt.Println("Tag created!")
}

func PushTags(repo *git.Repository, tagsToPush []string) {
	fmt.Printf("Pushing tags %+v...\n", tagsToPush)

	// Push tags in batches of BATCH_MAX_CAP, since pushing all together caused 
	// cloudbuild to not catch the push tag event for any of the new tags
	batches := make([][]config.RefSpec, 0)
	batch := make([]config.RefSpec, 0)
	batches = append(batches, batch)
	for _, tag := range tagsToPush {
		if len(batches[len(batches)-1]) == BATCH_MAX_CAP {
			batch := make([]config.RefSpec, 0)
			batches = append(batches, batch)
		}
		batches[len(batches)-1] = append(batches[len(batches)-1],
			config.RefSpec(fmt.Sprintf("refs/tags/%s:refs/tags/%s", tag, tag)))
	}

	for idx, batch := range batches {
		if idx > 0 {
			time.Sleep(time.Second * 2)
		}
		if err := repo.Push(&git.PushOptions{
			RemoteName: "origin",
			Progress:   os.Stdout,
			RefSpecs:   batch,
			Auth: &http.BasicAuth{
				Username: "ci-bot",
				Password: os.Getenv("GIT_ACCESS_TOKEN"),
			},
		}); err != nil {
			panic(err)
		}
	}
	fmt.Println("Release of tags completed!")
}

func CreateCommit(repo *git.Repository, message string, filepaths ...string) {
	fmt.Println("Creating commit...")
	tree, err := repo.Worktree()
	if err != nil {
		panic(err)
	}
	for _, file := range filepaths {
		if _, err := tree.Add(file); err != nil {
			panic(err)
		}
	}
	if _, err := tree.Commit(message, &git.CommitOptions{
		Author: &object.Signature{
			Name:  "Cloudbuild CI",
			Email: "technology@wopta.it",
			When:  time.Now(),
		},
	}); err != nil {
		panic(err)
	}
	fmt.Println("Commit created!")
}
