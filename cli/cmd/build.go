package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/ultimatesoftware/udocs/cli/storage"
	"github.com/ultimatesoftware/udocs/cli/udocs"
)

func Build() *cobra.Command {
	build := &cobra.Command{
		Use:   "build",
		Short: "Build a docs directory",
		Long: `
  udocs-build is for building a docs directory for local testing. It outputs rendered content in the
  directory '_docs'. README.md and SUMMARY.md files must exist in the root of the docs directory.
	`,
		Run: func(cmd *cobra.Command, args []string) {
			if err := udocs.Validate(dir); err != nil {
				fmt.Printf("Build failed: %v\n", err)
				return
			}

			os.RemoveAll("_docs")
			dao := storage.NewFileSystemDao("_docs", 0755, udocs.SearchPath())

			if err := udocs.Build(parseRoute(), dir, dao); err != nil {
				fmt.Printf("Build failed: %v\n", err)
				return
			}
		},
	}

	setFlag(build, "dir")
	return build
}
