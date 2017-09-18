package utils

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

//CHGRenameCmd defined the main command of RenameCHGenomicsFilename
var CHGRenameCmd = &cobra.Command{
	Use:   "RenameCHGenomicsFilename",
	Short: "A fastq.gz renamer from Cloudhealth Genomics",
	Long: `RenameCHGenomicsFilename
A fastq.gz renamer from Cloudhealth Genomics, avaliable command : 
	(1)rename;
	(2)rebuilddir;
Please use help for details

Author: Qinghui Li
Email: liqh@cloudhealth99.com
Company: CloudHealth Genomics`,

	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("No extra parameters allowed. Please re-check your input.")
			cmd.Help()
			return
		}
	},
}

func init() {
	CHGRenameCmd.AddCommand(versionCmd)
	CHGRenameCmd.AddCommand(rebuilddirCmd)
	CHGRenameCmd.AddCommand(renameCmd)
}

//Execute : start this program
func Execute() {
	if err := CHGRenameCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
