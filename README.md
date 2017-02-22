# RenameCHGenomicsFilename

### Author : Qinghui Li
- - -
### Decription
RenameCHGenomicsFilename is designed for rename fastq.gz files generated by CloudHealth Genomics.

The rules of fastq.gz files :
`[runID]_[machineId][flowcellID]_[CHGID]-[LibName]-[SampleName]-[Barcode]_[LaneId]_[ReadNum].fastq.gz`

It looks like this: S0001_01A_CHG000000-LIBNAME-SAMPLENAME-CCGGTTAA_L008_R2.fastq.gz

Some custumors do not want to get the fastq file in this style names. Some donot want the `CHGID` appears in the names, or `LibName`,`SampleName`,`Barcode`,`LaneId`.

So `RenameCHGenomicsFilename` is created.`RenameCHGenomicsFilename` is a `golang` based app.


###Usage
`RenameCHGenomicsFilename` can be run in both windows/linux OS, or anyother platforms which `golang` support.

##### Download
```git clone https://github.com/snailQH/RenameCHGenomicsFilename.git && cd RenameCHGenomicsFilename```

##### run the app from source code:
<pre><code>go run main.go -dir /online/projects/C170001-P001 -marker 5 //remove the LaneId from filename[/online/projects/C170001-P001]
</code></pre>

##### run the app in linux os:
<pre><code>./RenameCHGenomicsFilename -dir /online/projects/C170001-P001 -marker 5 //remove the LaneId from filename[/online/projects/C170001-P001]
</code></pre>

##### run the app in windows os:
<pre><code>RenameCHGenomicsFilename.exe -dir D:\\C170001-P001 -marker 5 //remove the LaneId from filename[D:\\C170001-P001]
</code></pre>

