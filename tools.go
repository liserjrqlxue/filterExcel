package main

import (
	"github.com/liserjrqlxue/goUtil/simpleUtil"
	"github.com/liserjrqlxue/goUtil/textUtil"
	"github.com/xuri/excelize/v2"
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
