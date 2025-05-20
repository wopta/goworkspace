package tag_functions

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"dev.wopta.it/cloudbuild/scripts/common"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

var (
	functions = flag.String("functions", "", "The functions to be released")
	fromEnv   = flag.String("from", "", "The source environment")
	targetEnv = flag.String("target", "", "The target environment")
)

func Exec() {
	flag.Parse()

	needsIncrement := false

	if *functions == "" || *targetEnv == "" {
		panic("The following arguments are required: --functions=<FUNCTIONS_TO_RELEASE> --target=<TARGET_ENV>")
	}
	if *fromEnv == "" {
		*fromEnv = *targetEnv
		needsIncrement = true
	}
	toUpdateFunctions := strings.Split(*functions, ",")

	repo, dir := common.CloneRepo()
	defer os.RemoveAll(dir)
	tagMap := fetchFunctionTags(repo, *fromEnv)

	tagsToPush := make([]string, 0)
	for _, fn := range toUpdateFunctions {
		currentTag := tagMap[fn][0]
		fmt.Printf("Latest tag for function %s - %+v\n", fn, currentTag)

		if !needsIncrement {
			checkoutTag(repo, currentTag)
		}

		newTag := generateNewTag(fn, currentTag.Version, needsIncrement)

		head, err := repo.Head()
		if err != nil {
			panic(err)
		}

		common.CreateTag(repo, newTag, fmt.Sprintf("Cloudbuild CI - Release %s in %s", fn, *targetEnv), head.Hash())

		tagsToPush = append(tagsToPush, newTag)
	}
	common.PushTags(repo, tagsToPush)
}

func checkoutTag(repo *git.Repository, tag common.GitTag) {
	fmt.Printf("Checking out tag %s ...\n", tag.Name)
	tree, err := repo.Worktree()
	if err != nil {
		panic(err)
	}
	if err := tree.Checkout(&git.CheckoutOptions{
		Branch: plumbing.ReferenceName(fmt.Sprintf("refs/tags/%s", tag.Name)),
	}); err != nil {
		panic(err)
	}
	fmt.Println("Checkout completed!")
}

func generateNewTag(module, baseTagVersion string, needsIncrement bool) string {
	fmt.Println("Generating new tag sequence...")
	baseTagArray := strings.Split(baseTagVersion, ".")
	patchNumber, err := strconv.ParseInt(baseTagArray[2], 10, 0)
	if err != nil {
		panic(err)
	}
	if needsIncrement {
		patchNumber++
	}
	tag := fmt.Sprintf("%s/%s.%s.%d.%s", module, baseTagArray[0], baseTagArray[1], patchNumber, *targetEnv)
	fmt.Printf("Generated sequence %s\n", tag)
	return tag
}

func fetchFunctionTags(repo *git.Repository, env string) map[string][]common.GitTag {
	fmt.Println("Fetching function tags ...")
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

		regex := regexp.MustCompile(fmt.Sprintf("([_a-zA-Z]*)/([[:digit:]].*).(%s)", env))
		if !regex.MatchString(t.Name) {
			return nil
		}

		tagParts := strings.Split(t.Name, "/")
		versionParts := strings.Split(tagParts[1], ".")

		tag := common.GitTag{
			Name:      t.Name,
			Module:    tagParts[0],
			Version:   strings.Join(versionParts[:len(versionParts)-1], "."),
			Env:       versionParts[len(versionParts)-1],
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
