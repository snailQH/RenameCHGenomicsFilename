package utils

import (
	"os"
	"strings"
	"time"
)

type logstype struct {
	Filename string
	Content  string
}

//WriteLogs write logs into logsfile
func (lg *logstype) WriteLogs() {
	fo, _ := os.Create(lg.Filename)
	defer fo.Close()
	fo.WriteString(lg.Content)
}

func (lg *logstype) getTime() {
	timestamp := time.Now().Unix()
	tm := time.Unix(timestamp, 0)
	x := tm.Format("2006-01-02  Mon 03:04:05 PM MST")
	lg.Content = x + "\n\n"
	x = tm.Format("2006-01-02 03:04:05")
	x = strings.Replace(x, " ", "", -1)
	x = strings.Replace(x, ":", "", -1)
	x = strings.Replace(x, "-", "", -1)
	lg.Filename = "RenameCHGenomicsFilename-" + x + ".log"
	//return logs, logsfile
}
