package tag_modules

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"dev.wopta.it/cloudbuild/scripts/common"
	"github.com/go-git/go-git/v5"
)

type GoModule struct {
	Name       string
	Version    string
	DependsOn  []string
	DependedBy []string
	Order      int
	IsModule   bool
	Dir        string
}

func (m *GoModule) UpdateSelf(repo *git.Repository) string {
	fmt.Printf("Current version %s\n", m.Version)
	versionParts := strings.Split(m.Version, ".")
	patchNumber, err := strconv.ParseInt(versionParts[2], 10, 0)
	if err != nil {
		panic(err)
	}
	newVersion := strings.Join([]string{
		versionParts[0],
		versionParts[1],
		fmt.Sprintf("%d", patchNumber+1),
	}, ".")
	newVersion = "v" + newVersion
	newTag := fmt.Sprintf("%s/%s", m.Name, newVersion)
	fmt.Printf("New version %s\n", newVersion)

	head, err := repo.Head()
	if err != nil {
		panic(err)
	}

	common.CreateTag(repo, newTag, fmt.Sprintf("Cloudbuild CI - Release module %s", newTag), head.Hash())

	m.Version = newVersion
	return newTag
}

func (m *GoModule) UpdateDependencies(repo *git.Repository) {
	for _, module := range m.DependedBy {
		goMod := common.GetGoMod(m.Dir, module)

		moduleImportPath := *modulePath + m.Name
		regex := regexp.MustCompile(fmt.Sprintf(
			"%s[[:space:]]v([[:digit:]]+.[[:digit:]]+.[[:digit:]]+)", moduleImportPath))

		fmt.Printf("Updating %s/go.mod with module %s - version %s\n", module, m.Name, m.Version)
		modGoMod := regex.ReplaceAllStringFunc(string(goMod), func(s string) string {
			return fmt.Sprintf("%s %s", moduleImportPath, m.Version)
		})

		if err := os.WriteFile(fmt.Sprintf("%s/%s/go.mod", m.Dir, module), []byte(modGoMod), os.ModePerm); err != nil {
			panic(err)
		}
		fmt.Println("go.mod updated!")

		common.CreateCommit(repo, fmt.Sprintf("chore: update %s deps", module), fmt.Sprintf("%s/go.mod", module))
	}
}
