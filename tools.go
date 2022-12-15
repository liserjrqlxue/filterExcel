package main

import (
	"github.com/liserjrqlxue/goUtil/simpleUtil"
	"github.com/xuri/excelize/v2"
	"strings"
)

func RemoveHitRow(excel *excelize.File, hitMap, includeDisease, excludeDisease map[string]bool, sheetName, hitCol, colName string) {
	var (
		hitIndex     int
		diseaseIndex int
	)
	var rows, err = excel.GetRows(sheetName)
	simpleUtil.CheckErr(err)

	for i, cell := range rows[0] {
		switch cell {
		case hitCol:
			hitIndex = i
		case colName:
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
				if includeDisease[disease] && excludeDisease[disease] {
					logInfo(i, hit, disease, "conflict!")
				} else if includeDisease[disease] {
					logInfo(i, hit, disease, "include")
					filter = false
				} else if excludeDisease[disease] {
					logInfo(i, hit, disease, "exclude")
				} else {
					logInfo(i, hit, disease, "lost!")
				}
			}
			if filter {
				logInfo(i, hit, diseaseInfo, "[remove]")
				simpleUtil.CheckErr(excel.RemoveRow(sheetName, i+1))
			} else {
				logInfo(i, hit, diseaseInfo, "[include]")
			}
		} else {
			logInfo(i, hit, diseaseInfo, "[noHit]")
		}
	}
}
