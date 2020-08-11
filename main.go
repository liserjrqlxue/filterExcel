package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/liserjrqlxue/goUtil/fmtUtil"
	"github.com/liserjrqlxue/goUtil/simpleUtil"
	"github.com/liserjrqlxue/goUtil/textUtil"
	"github.com/liserjrqlxue/version"
)

var (
	ex, _   = os.Executable()
	exPath  = filepath.Dir(ex)
	etcPath = filepath.Join(exPath, "etc")
)

var (
	input = flag.String(
		"input",
		"",
		"input excel",
	)
	output = flag.String(
		"output",
		"",
		"output excel, default -input.filter.xlsx",
	)
	hitList = flag.String(
		"list",
		"",
		"hit list to be filter",
	)
	sheetName = flag.String(
		"sheet",
		"All variants data",
		"sheet name to be filter",
	)
	colName = flag.String(
		"col",
		"疾病中文名",
		"column name of disease info",
	)
	sep = flag.String(
		"sep",
		"[n]",
		"sep to split disease info",
	)
	includeDisease = flag.String(
		"include",
		filepath.Join(etcPath, "includeDisease.list"),
		"include disease db file",
	)
	excludeDisease = flag.String(
		"exclude",
		filepath.Join(etcPath, "excludeDisease.list"),
		"exclude disease db file",
	)
	hitCol = flag.String(
		"hit",
		"SampleID",
		"column name of hit to be filter",
	)
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

	var hitIndex, diseaseIndex int
	var rows, err2 = inputExcel.GetRows(*sheetName)
	simpleUtil.CheckErr(err2)
	for i, cell := range rows[0] {
		switch cell {
		case *hitCol:
			hitIndex = i
		case *colName:
			diseaseIndex = i
		}
	}
	for i := len(rows) - 1; i > 0; i-- {
		var row = rows[i]
		var hit = row[hitIndex]
		var diseaseInfo = row[diseaseIndex]
		if hitMap[hit] {
			var diseaseInfos = strings.Split(diseaseInfo, *sep)
			var filter = true
			logInfo(i, hit, diseaseInfo, ":")
			for _, disease := range diseaseInfos {
				if includeDiseaseMap[disease] && excludeDiseaseMap[disease] {
					logInfo(i, hit, disease, "conflict!")
				} else if includeDiseaseMap[disease] {
					logInfo(i, hit, disease, "include")
					filter = false
				} else if excludeDiseaseMap[disease] {
					logInfo(i, hit, disease, "exclude")
				} else {
					logInfo(i, hit, disease, "lost!")
				}
			}
			if filter {
				logInfo(i, hit, diseaseInfo, "[remove]")
				simpleUtil.CheckErr(inputExcel.RemoveRow(*sheetName, i+1))
			} else {
				logInfo(i, hit, diseaseInfo, "[include]")
			}
		} else {
			logInfo(i, hit, diseaseInfo, "[noHit]")
		}
	}
	log.Printf("save as %s:%v", *output, inputExcel.SaveAs(*output))
	//simpleUtil.CheckErr(inputExcel.SaveAs(*output))
}

func logInfo(i int, sampleID, diseaseInfo, msg string) {
	log.Printf("row:%04d\t%s\t%s\t%s", i, sampleID, diseaseInfo, msg)
}
