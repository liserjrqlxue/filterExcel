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

	var cnvPackage, _ = textUtil.File2MapMap(filepath.Join(etcPath, "CNV包装.txt"), "产品编号", "\t", nil)

	var inputExcel, err = excelize.OpenFile(*input)
	simpleUtil.CheckErr(err)

	RemoveHitRows(inputExcel, *sheetName, *hitCol, *diseaseCol, *hitList, *includeDisease, *excludeDisease)
	// 补充实验
	MaskNotPackagedCNV(inputExcel, "补充实验", cnvPackage)

	log.Printf("save as %s:%v", *output, inputExcel.SaveAs(*output))
}

func logInfo(i int, sampleID, diseaseInfo, msg string) {
	log.Printf("row:%04d\t%s\t%s\t%s", i, sampleID, diseaseInfo, msg)
}
