package parsefilename

import (
	"io/ioutil"
	"os"
	"path"
	"strings"
)

func compileName(filename string, remove int, logs string) (string, string) {
	var result string
	dir := path.Dir(filename)
	name := path.Base(filename)
	rawstring := strings.Split(name, "_") //rawstring[0,1,2,3]SXXX _ XXB _ CHG000000-LIBNAME-SAMPLENAME-CCGGTTAA _ L00X _ R1.fastq.gz

	if len(rawstring) != 5 {
		logs = logs + "#Wrong samplename format :" + filename + "\n"
		return filename, logs
	}

	//remove chgid/libname/rawname/barcod/laneid[1/2/3/4/5]
	switch {
	case 1 <= remove && remove <= 3:
		secondstring := strings.Split(rawstring[2], "-")
		if len(secondstring) < 4 {
			logs = logs + "#Wrong samplename format :" + filename + "\n"
			return filename, logs
		}
		if strings.Contains(secondstring[0], "CHG") {
			secondstring = RemoveFromArray(secondstring, remove) //remove chg id
		} else {
			logs = logs + "#Wrong samplename: NO CHG ID IN THE RAW SAMPLENAME :" + filename + "\n"
		}
		rawstring[2] = strings.Join(secondstring, "-")
	case remove == 4:
		secondstring := strings.Split(rawstring[2], "-")
		if len(secondstring) < 4 {
			logs = logs + "#Wrong samplename format :" + filename + "\n"
			return filename, logs
		}
		if strings.Contains(secondstring[0], "CHG") {
			secondstring = RemoveFromArray(secondstring, len(secondstring)) //remove barcode
		} else {
			logs = logs + "#Wrong samplename: NO CHG ID IN THE RAW SAMPLENAME :" + filename + "\n"
		}
		rawstring[2] = strings.Join(secondstring, "-")
	case remove == 5:
		rawstring = RemoveFromArray(rawstring, 4)
	}
	result = path.Join(dir, strings.Join(rawstring, "_"))
	return result, logs
}

func writeLogs(logs string, logsfile string) {
	fo, _ := os.Create(logsfile)
	defer fo.Close()
	fo.WriteString(logs)
}

//ReName is design for compile the fastq filename
func ReName(filename []string, remove int) {
	checkfile := make(map[string]string)
	var logs string
	logsfile := "RenameCHGenomicsFilename.logs"

	for _, rawname := range filename {
		newname, logs := compileName(rawname, remove, logs)
		checkfile[rawname] = newname
		logs = logs + ""
	}

	for raw, new := range checkfile {
		if raw == new {
			logs = logs + "#Wrong format file\n"
			continue
		} else if checkfile[new] != "" {
			logs = logs + "#Duplicate filename in renaming " + raw + "\n"
			continue
		} else if checkfile[raw] == "renamed" {
			logs = logs + "#Duplicate filename in renaming " + raw + "\n"
			continue
		}
		logs = logs + ">#renaming " + raw + " > " + new + " ...\n"
		os.Rename(raw, new)
		checkfile[new] = "renamed"
	}
	writeLogs(logs, logsfile)
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
