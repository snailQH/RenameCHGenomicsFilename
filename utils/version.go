package utils

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of RenameCHGenomicsFilename",
	Long:  `All software has versions. This is RenameCHGenomicsFilename's`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(`
RenameCHGenomicsFilename Version
	RenameCHGenomicsFilename 2017-09-15
	details please refer to : https://github.com/snailQH/RenameCHGenomicsFilename`)
	},
}
