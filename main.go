package main

import (
	"flag"
	"github.com/liserjrqlxue/goUtil/fmtUtil"
	"github.com/liserjrqlxue/goUtil/simpleUtil"
	"github.com/liserjrqlxue/goUtil/textUtil"
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

	// load disease list
	var includeDiseases = textUtil.File2Array(*includeDisease)
	var includeDiseaseMap = make(map[string]bool)
	for _, d := range includeDiseases {
		includeDiseaseMap[d] = true
	}
	var excludeDiseases = textUtil.File2Array(*excludeDisease)
	var excludeDiseaseMap = make(map[string]bool)
	for _, d := range excludeDiseases {
		excludeDiseaseMap[d] = true
	}

	// load hit
	var hits = textUtil.File2Array(*hitList)
	var hitMap = make(map[string]bool)
	for _, h := range hits {
		hitMap[h] = true
	}

	var inputExcel, err1 = excelize.OpenFile(*input)
	simpleUtil.CheckErr(err1)
	RemoveHitRow(inputExcel, hitMap, includeDiseaseMap, excludeDiseaseMap, *sheetName, *hitCol, *checkCol)

	log.Printf("save as %s:%v", *output, inputExcel.SaveAs(*output))
}

func logInfo(i int, sampleID, diseaseInfo, msg string) {
	log.Printf("row:%04d\t%s\t%s\t%s", i, sampleID, diseaseInfo, msg)
}
