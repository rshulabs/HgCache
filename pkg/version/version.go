package version

import (
	"fmt"
	"os"

	"github.com/spf13/pflag"
)

const (
	FlagVersion          = "version"
	FlagVersionShorthand = "v"
)

var (
	GIT_TAG    string
	GIT_COMMIT string
	GIT_BRANCH string
	BUILD_TIME string
	GO_VERSION string
)

// FullVersion show the version info
func FullVersion() {
	fmt.Printf("Version   : %s\nBuild Time: %s\nGit Branch: %s\nGit Commit: %s\nGo Version: %s\n", GIT_TAG, BUILD_TIME, GIT_BRANCH, GIT_COMMIT, GO_VERSION)
	os.Exit(0)
}

// Short 版本缩写
func Short() {
	fmt.Printf("%s[%s %s]", GIT_TAG, BUILD_TIME, GIT_COMMIT)
	os.Exit(0)
}

func AddVersionFlag(name string, fs *pflag.FlagSet) *bool {
	return fs.BoolP(FlagVersion, FlagVersionShorthand, false, fmt.Sprintf("Version for %s.", name))
}
