package main

import (
	"flag"
	"github.com/liserjrqlxue/goUtil/fmtUtil"
	"github.com/liserjrqlxue/goUtil/simpleUtil"
	"github.com/liserjrqlxue/version"
	"github.com/xuri/excelize/v2"
	"log"
	"os"
	"path/filepath"
)

var (
	ex, _   = os.Executable()
	exPath  = filepath.Dir(ex)
	etcPath = filepath.Join(exPath, "etc")
)

func main() {
	version.LogVersion()
	flag.Parse()
	if *input == "" || *hitList == "" {
		flag.Usage()
		fmtUtil.Fprintln(os.Stderr, "-input and -list is required")
		os.Exit(1)
	}
	if *output == "" {
		*output = *input + ".filter.xlsx"
	}

	var inputExcel, err = excelize.OpenFile(*input)
	simpleUtil.CheckErr(err)

	RemoveHitRows(inputExcel, *sheetName, *hitCol, *diseaseCol, *hitList, *includeDisease, *excludeDisease)

	log.Printf("save as %s:%v", *output, inputExcel.SaveAs(*output))
}

func logInfo(i int, sampleID, diseaseInfo, msg string) {
	log.Printf("row:%04d\t%s\t%s\t%s", i, sampleID, diseaseInfo, msg)
}
