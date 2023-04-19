package main

import (
	"bufio"
	"flag"
	"log"
	"os"
	"path/filepath"
	"regexp"

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

	syrmList         = make(map[string]bool)
	syrmExcel        *excelize.File
	supplementReport = regexp.MustCompile(`补充报告`)
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

	for _, s := range textUtil.File2Array(filepath.Join(etcPath, "青岛崂山区9种单基因病携带者筛查项目.txt")) {
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
			} else if productID == "DX2063" && hospital == "十堰市人民医院" {
				syrmList[sampleID] = true
				log.Printf("%s\t%s\t%s\t十堰人民\n", sampleID, productID, hospital)
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

	if len(syrmList) > 0 {
		var (
			excelPath    = *output + ".syrm.xlsx"
			sheet1Name   = "All variants data"
			rows1        = simpleUtil.HandleError(inputExcel.GetRows(sheet1Name)).([][]string)
			title1       = rows1[0]
			sampleIndex1 = -1
			reportIndex  = -1
			rIdx         = 1

			sheet2Name   = "CNV"
			rows2        = simpleUtil.HandleError(inputExcel.GetRows(sheet2Name)).([][]string)
			title2       = rows2[0]
			sampleIndex2 = -1
		)

		for i, k := range title1 {
			if k == "SampleID" {
				sampleIndex1 = i
			}
			if k == "报告类别" {
				reportIndex = i
			}
		}

		for i, k := range title2 {
			if k == "#sample" {
				sampleIndex2 = i
			}
		}

		syrmExcel = excelize.NewFile()

		syrmExcel.NewSheet(sheet1Name)
		syrmExcel.NewSheet(sheet2Name)
		syrmExcel.DeleteSheet("Sheet1")

		writeRow(syrmExcel, sheet1Name, rows1[0], rIdx)
		rIdx++

		for i := len(rows1) - 1; i > 0; i-- {
			var (
				row        = rows1[i]
				sampleID   = row[sampleIndex1]
				reportType = row[reportIndex]
			)
			if syrmList[sampleID] && supplementReport.MatchString(reportType) {
				writeRow(syrmExcel, sheet1Name, row, rIdx)
				rIdx++
				simpleUtil.CheckErr(inputExcel.RemoveRow(sheet1Name, i+1))
			}
		}

		for i, row := range rows2 {
			var sampleID = row[sampleIndex2]
			if syrmList[sampleID] {
				writeRow(syrmExcel, sheet2Name, row, i+1)
			}
		}
		log.Printf("save as %s:%v", excelPath, syrmExcel.SaveAs(excelPath))
	}

	log.Printf("save as %s:%v", *output, inputExcel.SaveAs(*output))
}

func logInfo(i int, sampleID, diseaseInfo, msg string) {
	log.Printf("row:%04d\t%s\t%s\t%s", i, sampleID, diseaseInfo, msg)
}
