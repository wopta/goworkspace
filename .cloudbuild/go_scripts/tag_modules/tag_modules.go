package tag_modules

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"regexp"
	"slices"
	"strings"

	"dev.wopta.it/cloudbuild/scripts/common"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
)

var (
	modules    = flag.String("modules", "", "Comma separated list of modules to apply new tags")
	modulePath = flag.String("modulePath", "github.com/wopta/goworkspace/", "The path prefix for the go module")
)

func Exec() {
	flag.Parse()

	if *modules == "" {
		panic("Missing modules")
	}

	repo, dir := common.CloneRepo()
	defer os.RemoveAll(dir)
	tagMap := fetchModuleTags(repo)

	startingModules := strings.Split(*modules, ",")
	allModules := parseModules(dir, tagMap)
	modulesToUpdate := getDependentModules(allModules, startingModules)

	tagsToPush := make([]string, 0)
	for _, mod := range modulesToUpdate {
		fmt.Println("=========================================================")
		fmt.Printf("==== Working module %s...\n", mod.Name)

		tag := mod.UpdateSelf(repo)
		tagsToPush = append(tagsToPush, tag)

		mod.UpdateDependencies(repo)

		fmt.Printf("==== Module %s completed!\n", mod.Name)
		fmt.Println("=========================================================")
	}
	if err := repo.Push(&git.PushOptions{
		RemoteName: "origin",
		Progress:   os.Stdout,
		Auth: &http.BasicAuth{
			Username: "ci-bot",
			Password: os.Getenv("GIT_ACCESS_TOKEN"),
		},
	}); err != nil {
		panic(err)
	}
	common.PushTags(repo, tagsToPush)
}

func parseModules(path string, tagMap map[string][]common.GitTag) []GoModule {
	modulesToUpdate := make([]GoModule, 0)
	allModules := make([]string, 0)
	allModules = append(allModules, common.MODULES...)
	allModules = append(allModules, common.FUNCTIONS...)

	for _, m := range allModules {
		isModule := slices.Contains(common.MODULES, m)
		goMod := common.GetGoMod(path, m)

		regex := regexp.MustCompile(fmt.Sprintf("(?m)^[[:space:]]%s([_A-Za-z]*)[[:space:]](.*?)[[:space:]]", *modulePath))
		groups := regex.FindAllSubmatch(goMod, -1)
		module := GoModule{
			Name:     m,
			Order:    len(groups) - 1,
			IsModule: isModule,
			Dir:      path,
		}

		if tags, ok := tagMap[m]; ok {
			module.Version = tags[0].Version
		}

		dep := make([]string, 0)
		for _, matches := range groups {
			mod := string(matches[1])
			if slices.Contains(common.MODULES, mod) {
				dep = append(dep, mod)
			}
		}
		module.DependsOn = dep
		modulesToUpdate = append(modulesToUpdate, module)
	}

	for _, mod := range modulesToUpdate {
		for _, m := range mod.DependsOn {
			idx := slices.IndexFunc(modulesToUpdate, func(md GoModule) bool {
				return md.Name == m
			})
			modulesToUpdate[idx].DependedBy = append(modulesToUpdate[idx].DependedBy, mod.Name)
		}
	}

	slices.SortFunc(modulesToUpdate, func(a, b GoModule) int {
		return a.Order - b.Order
	})

	return modulesToUpdate
}

func getDependentModules(allModules []GoModule, startingModules []string) []GoModule {
	neededModules := make([]string, 0)
	for _, startingModule := range startingModules {
		moduleIndex := slices.IndexFunc(allModules, func(um GoModule) bool {
			return um.Name == startingModule
		})

		currentModule := allModules[moduleIndex]
		modulesToCheck := make([]string, 0)
		modulesToCheck = append(modulesToCheck, currentModule.Name)
		modulesToCheck = append(modulesToCheck, currentModule.DependedBy...)
		for _, mod := range modulesToCheck {
			if !slices.Contains(neededModules, mod) {
				neededModules = append(neededModules, mod)
			}
		}
	}

	return slices.DeleteFunc(allModules, func(m GoModule) bool {
		return !slices.Contains(neededModules, m.Name) || !m.IsModule
	})
}

func fetchModuleTags(repo *git.Repository) map[string][]common.GitTag {
	fmt.Println("Fetching module tags...")
	tagMap := make(map[string][]common.GitTag)

	iter, err := repo.Tags()
	if err != nil {
		panic(err)
	}
	if err := iter.ForEach(func(r *plumbing.Reference) error {
		t, err := repo.TagObject(r.Hash())
		if errors.Is(err, plumbing.ErrObjectNotFound) {
			return nil
		}
		if err != nil {
			panic(err)
		}

		regex := regexp.MustCompile("([_a-zA-Z]*)/v(([0-9.]*){3})")
		if !regex.MatchString(t.Name) {
			return nil
		}

		tagParts := strings.Split(t.Name, "/v")
		versionParts := strings.Split(tagParts[1], ".")

		tag := common.GitTag{
			Name:      t.Name,
			Module:    tagParts[0],
			Version:   strings.Join(versionParts, "."),
			Env:       "",
			Hash:      t.Hash.String(),
			CreatedAt: t.Tagger.When,
		}

		if _, ok := tagMap[tag.Module]; !ok {
			tagMap[tag.Module] = make([]common.GitTag, 0)
		}
		tagMap[tag.Module] = append(tagMap[tag.Module], tag)
		return nil
	}); err != nil {
		panic(err)
	}

	for _, tags := range tagMap {
		slices.SortStableFunc(tags, func(a, b common.GitTag) int {
			return b.CreatedAt.Compare(a.CreatedAt)
		})
	}
	fmt.Println("Fetch completed!")
	return tagMap
}
