package main

import (
	"flag"
	"fmt"

	"./parsefilename"
)

var marker int
var dir string

func init() {
	flag.IntVar(&marker, "marker", 4, "specify a marker to remove from filename[1/2/3/4/5]:chg id/original lib name/original samplename/barcode/lane id")
	flag.StringVar(&dir, "dir", "./", "specify a (relative/abs) path to fastq files")
}

func main() {
	flag.Usage = func() {
		fmt.Println("\nUsage of RenameCHGenomicsFilename:")
		fmt.Println("\nParameters:")
		fmt.Println("\n-marker :the region you want to remove from the samplenames[default:\"4\"]\n\t1: remove CHGID\n\t2: remove LibName\n\t3: remove SampleName\n\t4: remove Barcode\n\t5: remove LaneId")
		fmt.Printf("\n-dir :directory of the fastq files. eg: /online/projects/C150001-P001 [default:\"./\",for linux OS]\n\n")
		//flag.PrintDefaults()
	}
	flag.Parse()
	filelist := parsefilename.ListFile(dir)
	parsefilename.ReName(filelist, marker)
}
