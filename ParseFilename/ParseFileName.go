package ParseFilename

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

/**
SXXX_XXB_CHG000000-LIBNAME-SAMPLENAME-CCGGTTAA_L00X_R1.fastq.gz
SXXX_XXB_CHG000000-LIBNAME-SAMPLENAME-CCGGTTAA_L00X_R2.fastq.gz
SXXX_XXA_Undetermined_L00X_R1.fastq.gz;
SXXX_XXA_Undetermined_L00X_R2.fastq.gz
*/

func compileName(filename string, remove int) string {
	var result string
	dir := path.Dir(filename)
	name := path.Base(filename)
	rawstring := strings.Split(name, "_") //namestring[0,1,2,3]SXXX _ XXB _ CHG000000-LIBNAME-SAMPLENAME-CCGGTTAA _ L00X _ R1.fastq.gz

	if len(rawstring) != 5 {
		fmt.Printf("Wrong samplename format : %s.\n\tIt should be look like :SXXX_XXB_CHGXXXXXX-LIBNAME-SAMPLENAME-BARCODE_L00X_RX.fastq.gz\n", filename)
		return ""
	}

	//remove chgid/libname/rawname/barcod/laneid[1/2/3/4/5]
	switch {
	case 1 <= remove && remove <= 4:
		secondstring := strings.Split(rawstring[2], "-")
		if strings.Contains(secondstring[0], "CHG") {
			secondstring = RemoveFromArray(secondstring, remove) //remove chg id
		} else {
			fmt.Printf("Wrong samplename: NO CHG ID IN THE RAW SAMPLENAME %s !\n", filename)
		}
		rawstring[2] = strings.Join(secondstring, "-")
	case remove == 5:
		rawstring = RemoveFromArray(rawstring, 4)
	}
	result = path.Join(dir, strings.Join(rawstring, "_"))
	return result
}

//ReName is design for compile the fastq filename
func ReName(filename []string, remove int) {
	for _, rawname := range filename {
		fmt.Println(rawname)
		newname := compileName(rawname, remove)
		os.Rename(rawname, newname)
		//fmt.Println(">rename ", rawname, "to", newname)
	}
}

//ListFile is designed for list file in the dir
func ListFile(folder string) []string {
	var filelist []string
	files, _ := ioutil.ReadDir(folder)
	for _, file := range files {
		if file.IsDir() {
			newdir := path.Join(folder, file.Name())
			filelist = MergeArray(filelist, ListFile(newdir))
		} else {
			if strings.Contains(file.Name(), "Undetermined") {
				continue
			}
			x := path.Join(folder, file.Name())
			filelist = append(filelist, x)
		}
	}
	return filelist
}
