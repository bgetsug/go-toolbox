package stats

import (
	"fmt"

	"github.com/spf13/cobra"
)

var VersionCmd = &cobra.Command{
	Use:     "version",
	Aliases: []string{"build-metadata", "build-info"},
	Short:   "Display metadata about the current version/build",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("--------------------------------------")
		fmt.Println("            BUILD METADATA            ")
		fmt.Println("--------------------------------------")
		fmt.Printf("%-16s: %s\n", "Project Title", bmd.ProjectTitle)
		fmt.Printf("%-16s: %s\n", "Version", bmd.ProjectVersion)
		fmt.Printf("%-16s: %s\n", "VCS Branch", bmd.VCSBranch)
		fmt.Printf("%-16s: %s\n", "VCS Revision", bmd.VCSRevision)
		fmt.Printf("%-16s: %s\n", "Build Number", bmd.BuildNumber)
		fmt.Printf("%-16s: %s\n", "Build Timestamp", bmd.BuildTimestamp)
		fmt.Printf("%-16s: %s\n", "Environment", env)
		fmt.Println("--------------------------------------")
	},
}
