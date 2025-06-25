package lib

import (
	"fmt"
	"slices"
	"strconv"
	"strings"
)

func compareVersions(a, b string) int {
	version_a := strings.Split(a, ".")
	version_b := strings.Split(b, ".")

	if version_b[0] != version_a[0] {
		major_a, _ := strconv.ParseInt(version_a[0], 10, 0)
		major_b, _ := strconv.ParseInt(version_b[0], 10, 0)
		return int(major_b) - int(major_a)
	}
	if version_b[1] != version_a[1] {
		minor_a, _ := strconv.ParseInt(version_a[1], 10, 0)
		minnor_b, _ := strconv.ParseInt(version_b[1], 10, 0)
		return int(minnor_b) - int(minor_a)
	}

	patch_a_raw, _, _ := strings.Cut(version_a[2], "-")
	patch_b_raw, _, _ := strings.Cut(version_b[2], "-")
	patch_a, _ := strconv.ParseInt(patch_a_raw, 10, 0)
	patch_b, _ := strconv.ParseInt(patch_b_raw, 10, 0)
	return int(patch_b) - int(patch_a)
}

func GetLastVersionPrecontrattuale(productName, version string) (string, error) {
	bucket := "documents-public-dev"
	path := fmt.Sprint("information-sets/", productName, "/", version)
	res, err := getLastVersionDocument(bucket, path, "Precontrattuale")
	if err != nil {
		return "", err
	}
	if res == "none" { //return default pdf
		return fmt.Sprint(bucket, "/", path, "/Precontrattuale.pdf"), err
	}
	return res, nil
}

func getLastVersionDocument(bucket, rootPath, versioningDirectory string) (string, error) {
	rootPath += "/" + versioningDirectory + "/"
	fileList, err := ListGoogleStorageFolderContentWithBucket(rootPath, bucket)
	fileList = slices.DeleteFunc(fileList, func(path string) bool {
		return path == rootPath
	})
	if err != nil {
		if strings.Contains(err.Error(), "no such file or directory") {
			return "none", nil
		}
		return "", err
	}
	if len(fileList) == 0 {
		return "none", nil
	}
	for i := range fileList {
		parts := strings.SplitAfter(fileList[i], "/")
		fileList[i] = parts[len(parts)-1]
	}
	slices.SortFunc(fileList, compareVersions)
	return fmt.Sprint(bucket, "/", rootPath, fileList[0]), nil
}
