package utils

import (
	"fmt"
	"os"
)

//RemoveFromArray is designed for remove one item from a string array
func RemoveFromArray(array []string, index int) []string {
	var newarray []string
	//index : 1~len(array)
	//SXXX _ XXB _ CHG000000-LIBNAME-SAMPLENAME-CCGGTTAA _ L00X _ R1.fastq.gz
	if index <= 0 {
		fmt.Println("out of range")
		os.Exit(1)
	} else if index > len(array) {
		fmt.Println("out of range")
		os.Exit(1)
	} else if index == 1 {
		newarray = array[1:]
	} else if index == len(array) {
		newarray = array[:(len(array) - 1)]
	} else {
		a := array[0:(index - 1)]
		b := array[index:]
		newarray = append(a, b...)
	}
	return newarray
}
