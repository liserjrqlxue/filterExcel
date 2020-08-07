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
		"Disease*",
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
	for i, row := range rows {
		if i == 0 {
			continue
		}
		var hit = row[hitIndex]
		var diseaseInfo = row[diseaseIndex]
		if hitMap[hit] {
			var diseaseInfos = strings.Split(diseaseInfo, *sep)
			var filter = true
			//log.Printf("%s\t%s:",hit,diseaseInfo)
			for _, disease := range diseaseInfos {
				if includeDiseaseMap[disease] && excludeDiseaseMap[disease] {
					log.Printf("%-12s\t%s\tconflict!", hit, disease)
				} else if includeDiseaseMap[disease] {
					//log.Printf("%s\t%s\tinclude",hit,disease)
					filter = false
				} else if excludeDiseaseMap[disease] {
					//log.Printf("%s\t%s\texclude",hit,disease)
				} else {
					log.Printf("%-12s\t%s\tlost!", hit, disease)
				}
			}
			if filter {
				//log.Printf("%s\t%s\t%d\t[remove]",hit,diseaseInfo,i)
				simpleUtil.CheckErr(inputExcel.RemoveRow(*sheetName, i))
			} else {
				//log.Printf("%s\t%s\t%d\t[include]",hit,diseaseInfo,i)
			}
		} else {
			//log.Printf("%s\t%s\t%d\t[noHit]",hit,diseaseInfo,i)
		}
	}
	simpleUtil.CheckErr(inputExcel.SaveAs(*output))
}
