package main

import (
	"github.com/liserjrqlxue/goUtil/simpleUtil"
	"github.com/liserjrqlxue/goUtil/textUtil"
	"github.com/xuri/excelize/v2"
	"regexp"
	"strings"
)

func RemoveHitRows(excel *excelize.File, sheetName, hitCol, diseaseCol, hitList, includeList, excludeList string) {
	var (
		hitIndex     int
		diseaseIndex int
		hitMap       = List2BoolMap(hitList)
		include      = List2BoolMap(includeList)
		exclude      = List2BoolMap(excludeList)
	)
	var rows, err = excel.GetRows(sheetName)
	simpleUtil.CheckErr(err)

	for i, cell := range rows[0] {
		switch cell {
		case hitCol:
			hitIndex = i
		case diseaseCol:
			diseaseIndex = i
		}
	}

	for i := len(rows) - 1; i > 0; i-- {
		var (
			row         = rows[i]
			hit         = row[hitIndex]
			diseaseInfo = row[diseaseIndex]
		)
		if hitMap[hit] {
			RemoveHitRow(excel, sheetName, diseaseInfo, *sep, hit, i, include, exclude)
		}
	}
}

func RemoveHitRow(excel *excelize.File, sheetName, diseases, sep, hit string, i int, include, exclude map[string]bool) {
	var filter = true
	logInfo(i, hit, diseases, ":")
	for _, disease := range strings.Split(diseases, sep) {
		if include[disease] && exclude[disease] {
			logInfo(i, hit, disease, "conflict!")
		} else if include[disease] {
			logInfo(i, hit, disease, "include")
			filter = false
		} else if exclude[disease] {
			logInfo(i, hit, disease, "exclude")
		} else {
			logInfo(i, hit, disease, "lost!")
		}
	}
	if filter {
		logInfo(i, hit, diseases, "[remove]")
		simpleUtil.CheckErr(excel.RemoveRow(sheetName, i+1))
	} else {
		logInfo(i, hit, diseases, "[include]")
	}
}

func List2BoolMap(path string) map[string]bool {
	var boolMap = make(map[string]bool)
	for _, s := range textUtil.File2Array(path) {
		boolMap[s] = true
	}
	return boolMap
}

var (
	isThal = regexp.MustCompile(`地贫`)
	isSMA  = regexp.MustCompile(`SMN1`)
	isF8   = regexp.MustCompile(`F8`)
)

func MaskNotPackagedCNV(excel *excelize.File, sheetName string, packages map[string]map[string]string, smaList map[string]bool) {
	var (
		sampleIndex int
		hitIndex    int
		thalIndexs  []int
		smaIndexs   []int
		f8Indexs    []int
		maskValue   = "检测范围外"
		rows, err   = excel.GetRows(sheetName)
	)
	simpleUtil.CheckErr(err)

	for i, cell := range rows[0] {
		if cell == "SampleID" {
			sampleIndex = i
		}
		if cell == "产品编码_产品名称" {
			hitIndex = i
		}
		if isThal.MatchString(cell) {
			thalIndexs = append(thalIndexs, i)
		}
		if isSMA.MatchString(cell) {
			smaIndexs = append(smaIndexs, i)
		}
		if isF8.MatchString(cell) {
			f8Indexs = append(f8Indexs, i)
		}
	}
	for i, row := range rows {
		if i == 0 {
			continue
		}
		var hit = strings.Split(row[hitIndex], "_")[0]
		var info = packages[hit]
		var sampleID = row[sampleIndex]
		if info["地贫"] != "是" {
			maskCells(excel, sheetName, maskValue, i, thalIndexs)
		}
		if info["SMA"] != "是" || smaList[sampleID] {
			maskCells(excel, sheetName, maskValue, i, smaIndexs)
		}
		if info["F8"] != "是" {
			maskCells(excel, sheetName, maskValue, i, f8Indexs)
		}
	}
}

func maskCells(excel *excelize.File, sheetName, maskValue string, rIdx int, cols []int) {
	for _, col := range cols {
		simpleUtil.CheckErr(
			excel.SetCellStr(sheetName, GetAxis(col+1, rIdx+1), maskValue),
		)
	}
}

func GetAxis(col, row int) string {
	var axis, err = excelize.CoordinatesToCellName(col, row)
	simpleUtil.CheckErr(err)
	return axis
}
