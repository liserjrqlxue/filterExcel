package main

import (
	"bufio"
	"flag"
	"log"
	"os"
	"path/filepath"

	"github.com/liserjrqlxue/goUtil/fmtUtil"
	"github.com/liserjrqlxue/goUtil/osUtil"
	"github.com/liserjrqlxue/goUtil/scannerUtil"
	"github.com/liserjrqlxue/goUtil/simpleUtil"
	"github.com/liserjrqlxue/goUtil/textUtil"
	"github.com/liserjrqlxue/version"
	"github.com/xuri/excelize/v2"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

var (
	ex, _   = os.Executable()
	exPath  = filepath.Dir(ex)
	etcPath = filepath.Join(exPath, "etc")
)

var (
	lsmsList        = make(map[string]bool)
	lsmsHopitalList = make(map[string]bool)
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

	for _, s := range textUtil.File2Array(filepath.Join(etcPath, "")) {
		lsmsHopitalList[s] = true
	}
	if *addition != "" {
		var enc = simplifiedchinese.GBK
		var file = osUtil.Open(*addition)
		var r = transform.NewReader(file, enc.NewDecoder())
		var scanner = bufio.NewScanner(r)
		for _, strings := range scannerUtil.Scanner2Slice(scanner, "\t") {
			var sampleID = strings[0]
			var productID = strings[11]
			var hospital = strings[14]
			if productID == "DX2063" && lsmsHopitalList[hospital] {
				lsmsList[sampleID] = true
				log.Printf("%s\t%s\t%s\t崂山民生\n", sampleID, productID, hospital)
			} else {
				log.Printf("%s\t%s\t%s\n", sampleID, productID, hospital)
			}
		}
		simpleUtil.CheckErr(file.Close())
	}

	if *lsms != "" {
		for _, s := range textUtil.File2Array(*lsms) {
			lsmsList[s] = true
		}
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
