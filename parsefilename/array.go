package parsefilename

//RemoveFromArray is designed for remove one item from a string array
func RemoveFromArray(array []string, index int) []string {
	var newarray []string
	//index : 1~len(array)
	if index <= 0 {
		panic("out of range")
	} else if index > len(array) {
		panic("out of range")
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
