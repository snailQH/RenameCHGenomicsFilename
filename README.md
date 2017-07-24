# RenameCHGenomicsFilename

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

```-dir``` : set the directory of the fastq files

```-marker``` : specify the element you want to remove from the filename

&#8195;&#8195;&#8195;&#8195;0: remove RunID, flowcellID, CHGID, and barcode Info
&#8195;&#8195;&#8195;&#8195;1: remove CHGID
&#8195;&#8195;&#8195;&#8195;2: remove LibName
&#8195;&#8195;&#8195;&#8195;3: remove SampleName
&#8195;&#8195;&#8195;&#8195;4: remove Barcode
&#8195;&#8195;&#8195;&#8195;5: remove LaneId

##### Download
```git clone https://github.com/snailQH/RenameCHGenomicsFilename.git && cd RenameCHGenomicsFilename```

##### 1. run the app from source code:
<pre><code>go run main.go -dir /online/projects/C170001-P001 -marker 5 //remove the LaneId from filename[/online/projects/C170001-P001]
</code></pre>

##### 2. run the app in linux os:
<pre><code>./RenameCHGenomicsFilename -dir /online/projects/C170001-P001 -marker 5 //remove the LaneId from filename[/online/projects/C170001-P001]
</code></pre>

##### 3. run the app in windows os:
<pre><code>RenameCHGenomicsFilename.exe -dir D:\\C170001-P001 -marker 5 //remove the LaneId from filename[D:\\C170001-P001]
</code></pre>

##### 4. run the app directly:
In this way, you can remove the barcode from all the `*.fastq.gz`files in the current and sub directory.
<pre><code>RenameCHGenomicsFilename.exe	//in windows os
./RenameCHGenomicsFilename	//in linux os
</code></pre>


