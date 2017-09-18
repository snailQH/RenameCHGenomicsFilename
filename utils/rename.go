package utils

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
)

var marker int
var dir string

var renameCmd = &cobra.Command{
	Use:   "rename",
	Short: "rename is used for rename fastq file(s) to user specified stype",
	Long: `rename is used for rename fastq file(s) to user specified stype, avaliable Parameters: 
	(1)-marker :the region you want to remove from the samplenames[default:"4"]
		0: remove RunId,flowcellID,CHGID,Barcode
		1: remove CHGID
		2: remove LibName
		3: remove SampleName
		4: remove Barcode
		5: remove LaneId
		12: remove CHGID & LibName
		14: remove CHGID & Barcode
		124: remove CHGID & LibName & Barcode
	(2)-dir :directory of the fastq files. eg: /online/projects/C150001-P001 [default:"./",means current dir,for linux OS]
	Any problem please contact  Qinghui Li  via liqh@cloudhealth99.com
	
	`,

	Run: func(cmd *cobra.Command, args []string) {
		//for rename module will do the action of rename file in the current dir, so no extra command checking
		//no extra parameters allowed,check the marker value,whether the dir exist
		if len(args) != 0 || !PathExist(dir) || !(marker == 1 || marker == 2 || marker == 3 || marker == 4 || marker == 5 || marker == 12 || marker == 14 || marker == 124) {
			if len(args) != 0 {
				fmt.Println("#No extra parameters allowed. Please re-check your input.")
			}

			if !PathExist(dir) {
				fmt.Printf("#Cannot find dir: %s. Please re-check your input.\n", dir)
			}

			if !(marker == 1 || marker == 2 || marker == 3 || marker == 4 || marker == 5 || marker == 12 || marker == 14 || marker == 124) {
				fmt.Printf("#Wrong marker input :  %d. Please re-check your input.\n", marker)
			}

			cmd.Help()
			return
		}
		ReName(ListFile(dir), marker) //do rename file here
	},
}

func init() {
	renameCmd.Flags().IntVarP(&marker, "marker", "m", 4, "specify a marker to remove from filename[0/1/2/3/4/5/12/14/124]:chg id/original lib name/original samplename/barcode/lane id/CHGID & LibName/CHGID & Barcode/CHGID & LibName & Barcode")
	renameCmd.Flags().StringVarP(&dir, "dir", "d", "./", "specify a (relative/abs) path to fastq files")
}

func compileName(filename string, remove int, logs string) (string, string) {
	var result string
	dir := path.Dir(filename)
	name := path.Base(filename)
	rawstring := strings.Split(name, "_") //rawstring[0,1,2,3,4] SXXX _ XXB _ CHG000000-LIBNAME-SAMPLENAME-CCGGTTAA _ L00X _ R1.fastq.gz

	if len(rawstring) != 5 {
		logs = logs + "#Wrong samplename format: " + filename + "\n"
		return filename, logs
	}

	//remove chgid/libname/rawname/barcod/laneid[1/2/3/4/5]

	secondstring := strings.Split(rawstring[2], "-") //CHG000000-LIBNAME-SAMPLENAME-CCGGTTAA
	if len(secondstring) < 4 {
		logs = logs + "#Wrong samplename format: " + filename + "\n"
		return filename, logs
	}

	switch {
	case 1 <= remove && remove <= 3:
		if strings.Contains(secondstring[0], "CHG") {
			secondstring = RemoveFromArray(secondstring, remove) //remove chg id
		} else {
			logs = logs + "#Wrong samplename: NO CHG ID IN THE RAW SAMPLENAME: " + filename + " \n"
		}
		rawstring[2] = strings.Join(secondstring, "-")
	case remove == 4:
		if strings.Contains(secondstring[0], "CHG") {
			secondstring = RemoveFromArray(secondstring, len(secondstring)) //remove barcode
			if dualBarcode(secondstring[len(secondstring)-1]) {
				secondstring = RemoveFromArray(secondstring, len(secondstring)) //remove dual barcode
			}
		} else {
			logs = logs + "#Wrong samplename: NO CHG ID IN THE RAW SAMPLENAME: " + filename + " \n"
		}
		rawstring[2] = strings.Join(secondstring, "-")
	case remove == 5:
		rawstring = RemoveFromArray(rawstring, 4) //remove LaneID

	case remove == 0:
		//remove all cloudhealth genomics sampleinfo
		newstring := secondstring[1:(len(secondstring) - 1)] //remove CHG ID and barcode info
		new := strings.Join(newstring, "-")
		result := new + "_" + rawstring[3] + "_" + rawstring[4]
		result = path.Join(dir, result) //rename the file in the original dir
		return result, logs

	case remove == 12:
		secondstring = RemoveFromArray(secondstring, 1) //remove CHG ID
		secondstring = RemoveFromArray(secondstring, 1) //remove Lib ID
		rawstring[2] = strings.Join(secondstring, "-")
	case remove == 14:
		secondstring = RemoveFromArray(secondstring, 1)                 //remove CHG ID
		secondstring = RemoveFromArray(secondstring, len(secondstring)) //remove barcode
		if dualBarcode(secondstring[len(secondstring)-1]) {
			secondstring = RemoveFromArray(secondstring, len(secondstring)) //remove dual barcode
		}
		rawstring[2] = strings.Join(secondstring, "-")
	case remove == 124:
		secondstring = RemoveFromArray(secondstring, 1)                 //remove CHG ID
		secondstring = RemoveFromArray(secondstring, 1)                 //remove Lib ID
		secondstring = RemoveFromArray(secondstring, len(secondstring)) //remove barcode
		if dualBarcode(secondstring[len(secondstring)-1]) {
			secondstring = RemoveFromArray(secondstring, len(secondstring)) //remove dual barcode
		}
		rawstring[2] = strings.Join(secondstring, "-")
	}
	result = path.Join(dir, strings.Join(rawstring, "_"))
	return result, logs
}

func dualBarcode(text string) bool {
	barcoderegexp := regexp.MustCompile("[ATGC]+")
	result := barcoderegexp.FindString(text)
	if result != "" {
		return true
	}
	return false
}

//ReName is design for compile the fastq filename
func ReName(filename []string, remove int) {
	checkfile := make(map[string]string)
	var renameLogs logstype
	renameLogs.getTime()
	var logs = renameLogs.Content
	//var logs, logsfile string
	//logs, logsfile = getTime(logs, logsfile)

	for _, rawname := range filename {
		newname, logs := compileName(rawname, remove, logs)
		checkfile[rawname] = newname
		renameLogs.Content = logs + ""
	}

	for raw, new := range checkfile {
		if raw == new {
			renameLogs.Content = renameLogs.Content + "#Wrong format file: " + raw + "\n"
			continue
		} else if checkfile[new] != "" {
			renameLogs.Content = renameLogs.Content + "#Duplicate filename in renaming: " + raw + "\n"
			continue
		} else if checkfile[raw] == "renamed" {
			renameLogs.Content = renameLogs.Content + "#Duplicate filename in renaming: " + raw + "\n"
			continue
		}
		renameLogs.Content = renameLogs.Content + ">#renaming " + raw + " > " + new + " ...\n"
		os.Rename(raw, new)
		checkfile[new] = "renamed"
	}
	renameLogs.WriteLogs()
}

//ListFile is designed for list file in the dir
func ListFile(folder string) []string {
	var filelist []string
	files, _ := ioutil.ReadDir(folder)
	for _, file := range files {
		if file.IsDir() {
			newdir := path.Join(folder, file.Name())
			filelist = append(filelist, ListFile(newdir)...)
		} else {
			if strings.Contains(file.Name(), "Undetermined") {
				continue
			} else if strings.Contains(file.Name(), "fastq.gz") {
				x := path.Join(folder, file.Name())
				filelist = append(filelist, x)
			}
		}
	}
	return filelist
}

//PathExist is for determine whether the file is exist
func PathExist(path string) bool {
	_, err := os.Stat(path)
	if err != nil && os.IsNotExist(err) {
		return false
	}
	return true
}
