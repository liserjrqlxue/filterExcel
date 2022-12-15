package main

import (
	"flag"
	"path/filepath"
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
	checkCol = flag.String(
		"col",
		"疾病中文名",
		"column name to disease to check filter",
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
