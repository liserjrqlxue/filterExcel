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
	if *input == "" {
		flag.Usage()
		fmtUtil.Fprintln(os.Stderr, "-input is required")
		os.Exit(1)
	}
	if *output == "" {
		*output = *input + ".filter.xlsx"
	}

	var lsmsList = make(map[string]bool)
	for _, s := range textUtil.File2Array(*lsms) {
		lsmsList[s] = true
	}

	var cnvPackage, _ = textUtil.File2MapMap(filepath.Join(etcPath, "CNV包装.txt"), "产品编号", "\t", nil)

	var inputExcel, err = excelize.OpenFile(*input)
	simpleUtil.CheckErr(err)

	if *hitList != "" {
		RemoveHitRows(inputExcel, *sheetName, *hitCol, *diseaseCol, *hitList, *includeDisease, *excludeDisease)
	}
	// 补充实验
	MaskNotPackagedCNV(inputExcel, "补充实验", cnvPackage, lsmsList)

	log.Printf("save as %s:%v", *output, inputExcel.SaveAs(*output))
}

func logInfo(i int, sampleID, diseaseInfo, msg string) {
	log.Printf("row:%04d\t%s\t%s\t%s", i, sampleID, diseaseInfo, msg)
}
