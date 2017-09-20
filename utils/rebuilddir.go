package utils

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

type rebuildtype struct {
	ProjectID  string //cur projectid to be rebuild
	FileDir    string //the dir of the disk
	CustomerID string //customer id,means the re-build style
}

var customerid, projectid, projectdir string
var projectstyle = regexp.MustCompile("[CVRUHX]\\d{6}")           //C150001
var samplestyle = regexp.MustCompile("CHG\\d{6}")                 //CHG020000
var filetype = regexp.MustCompile("S\\d{4}_\\d{2}[AB]_CHG\\d{6}") //S0618_01B_CHG026244-WHTRDRPEP00004838-WHRDMETmgpMAAAAAA-66-AGGTTAAC_L003_R1.fastq.gz

var rebuilddirCmd = &cobra.Command{
	Use:   "rebuild",
	Short: "re-build the structure of the delivering fastq files",
	Long:  "re-build the structure of the delivering fastq files，including fq files,md5 files, dir name, et.al.",

	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 0 || !PathExist(projectdir) {
			if len(args) != 0 {
				fmt.Println("#No extra parameters allowed. Please re-check your input.")
			}
			if !IsRightCustomerID(customerid) {
				fmt.Println("#Wrong CustomerID input. Please re-check your input.")
			}
			if !IsRightProjectID(projectid) {
				fmt.Println("#Wrong ProjectID input. Please re-check your input.")
			}
			if !PathExist(projectdir) {
				fmt.Printf("#Cannot find projectdir: %s. Please re-check your input.\n", projectdir)
			}
			cmd.Help()
			return
		}
		var rb rebuildtype
		rb.FileDir = projectdir
		rb.ProjectID = projectid
		rb.CustomerID = customerid
		rb.ReBuild()
	},
}

func init() {
	rebuilddirCmd.Flags().StringVarP(&customerid, "customerid", "c", "", `customer id , determined the style of fastq files structure. 
		Please input the correct customerid, like: C150001,C150003`)
	rebuilddirCmd.Flags().StringVarP(&projectdir, "dir", "d", "", `specify a (relative/abs) path to projectdir,like C150003-P999`)
	rebuilddirCmd.Flags().StringVarP(&projectid, "projectid", "p", "all", `projectid of the current delieverying`) //default rebuild all the project-dir in the specified dir
}

//IsRightCustomerID make a judgement for the input CustomerId
func IsRightCustomerID(customerid string) bool {
	var style = regexp.MustCompile("[CVRUHX]\\d{6}") //C150001-P001
	if style.MatchString(customerid) {
		return true
	}
	return false
}

//IsRightProjectID make a judgement for the input projectid
func IsRightProjectID(projectid string) bool {
	var style = regexp.MustCompile("[CVRUHX]\\d{6}-P\\d{3}") //C150001-P001
	if style.MatchString(projectid) {
		return true
	}
	return false
}

//FindSubDir is designed for list file/dir in the dir,dir bool define whether dir or file
func FindSubDir(folder string, style *regexp.Regexp, dir bool) []string {
	var subdirs []string
	files, _ := ioutil.ReadDir(folder)
	for _, file := range files {
		if dir {
			if file.IsDir() {
				if style.MatchString(file.Name()) {
					subdirs = append(subdirs, file.Name())
				}
			}
		} else {
			if file.IsDir() {
				continue
			} else {
				if style.MatchString(file.Name()) {
					subdirs = append(subdirs, file.Name())
				}
			}
		}
	}
	return subdirs
}

//GetCurDate GetCurDate like 20170916
func GetCurDate() string {
	timestamp := time.Now().Unix()
	tm := time.Unix(timestamp, 0)
	fmt.Println(tm)
	x := tm.Format("2006-01-02")
	x = strings.Replace(x, "-", "", -1)
	return x
}

//ReBuild will re-build the dir and file according to the customer's rules
func (rb *rebuildtype) ReBuild() {
	var rebuildlogs logstype
	rebuildlogs.getTime()
	var logs = rebuildlogs.Content
	logs = logs + fmt.Sprintf("\n##do the action of re-build file for %s %s\n", rb.ProjectID, rb.FileDir)
	allprojectdir := FindSubDir(rb.FileDir, projectstyle, true) //get all the project dir in the input dir:rb.FileDir

	for _, rawprojectdir := range allprojectdir {
		if projectid != "all" {
			if !strings.Contains(rawprojectdir, projectid) {
				continue
			}
		}
		logs = rebuildbycustomid(rb.CustomerID, rawprojectdir, logs)
	}
	rebuildlogs.WriteLogs()
}

func findsublibname(filelist []string) string {
	var libname, tmpname = "", "" //tmpname will record the current right string match all the samplenames
	var subpattern []string
	for _, file := range filelist {
		if subpattern == nil && strings.Contains(file, "-") {
			subpattern = strings.Split(filelist[0], "-")
			break
		}
	}

	//var subpattern = strings.Split(filelist[0], "-")
	for a, b := range subpattern {
		if a > 0 {
			if libname == "" {
				tmpname = subpattern[1]
			} else {
				tmpname = libname + "-" + b
			}
			flag := checklibname(filelist, tmpname, a)
			if flag {
				libname = tmpname
			} else {
				if a == 1 {
					fmt.Println("#error:cannot find the sublibname")
				}
				return libname
			}
		}
	}
	return libname
}

func checklibname(filelist []string, libname string, index int) bool {
	var flag = true
	for _, fastqfilename := range filelist {
		var subpattern = strings.Split(fastqfilename, "-")
		tmparray := subpattern[1:(index + 1)]
		curname := strings.Join(tmparray, "-")
		if libname != curname {
			return false // find the unmatched pattern in current samplename
		}
	}
	return flag
}

func rebuildbycustomid(customerid string, rawprojectdir string, logs string) string {
	if customerid == "C150003" {
		newprojectdir := rawprojectdir + "-" + GetCurDate()
		if PathExist(newprojectdir) {
			if err := os.RemoveAll(newprojectdir); err != nil {
				logs = logs + fmt.Sprintf("\n#error in remove exist newprojectdir %v", err)
			}
		}

		os.Mkdir(newprojectdir, 0774) //mkdir for the new dirname,project level
		rawdatadir := path.Join(newprojectdir, "Rawdata")
		os.Mkdir(rawdatadir, 0774)

		_, projectname := path.Split(rawprojectdir)
		filemd5 := rawprojectdir + ".local.md5"
		filecheckmd5 := rawprojectdir + "local.md5.check"
		if !PathExist(filemd5) || !PathExist(filecheckmd5) {
			logs = logs + fmt.Sprintf("\n#error,md5 file or md5.check file missing")
		} else {
			newfilemd5 := projectname + ".local.md5"
			newfilecheckmd5 := projectname + "local.md5.check"
			newfilemd5 = path.Join(newprojectdir, filemd5)
			newfilecheckmd5 = path.Join(newprojectdir, filecheckmd5)
			os.Rename(filemd5, newfilemd5)
			os.Rename(filecheckmd5, newfilecheckmd5)
		}

		allsampledir := FindSubDir(rawprojectdir, samplestyle, true) //find sample dir in the raw project dir
		for _, sampledir := range allsampledir {
			sampledir = path.Join(rawprojectdir, sampledir)
			filelist := FindSubDir(sampledir, filetype, false)
			libname := findsublibname(filelist)
			logs = logs + fmt.Sprintf("\n##sampledir: %s -> matched libname : %s\n", sampledir, libname)
			libdir := path.Join(rawdatadir, libname)
			os.Mkdir(libdir, 0774)
			for _, file := range filelist {
				newfile := path.Base(file)
				newfile = path.Join(libdir, newfile)
				logs = logs + fmt.Sprintf("##Renaming file: %s %s \n", file, newfile)
				if err := os.Rename(file, newfile); err != nil {
					logs = logs + fmt.Sprintf("#error in link rawfile to newfile %s \n", err)
				}
			}
		}
		if err := os.RemoveAll(rawprojectdir); err != nil {
			logs = logs + fmt.Sprintf("#error in remove raw project dir %s\n", err)
		}
	}

	return logs
}
