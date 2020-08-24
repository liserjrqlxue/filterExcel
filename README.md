# filterExcel
~~过滤掉附件表格中的55种疾病~~  
对于指定样品保留附件表格中的100种疾病

## build
```shell script
gitDescribe=$(git branch --show-current):$(git describe --tags)
golangVersion=$(go version)
buildStamp=$(date -u '+%Y-%m-%d_%I:%M:%S%p')
go build -x -ldflags "-s -w -X 'github.com/liserjrqlxue/version.gitDescribe=$gitDescribe' -X 'github.com/liserjrqlxue/version.buildStamp=$buildStamp' -X 'github.com/liserjrqlxue/version.golangVersion=$golangVersion'"
```

## usage
```
Usage of /zfsyt1/B2C_RD_P2/USER/wangyaoshen/pipeline/filterExcel/filterExcel:
  -col string
    	column name of disease info (default "疾病中文名")
  -exclude string
    	exclude disease db file (default "/zfsyt1/B2C_RD_P2/USER/wangyaoshen/pipeline/filterExcel/etc/excludeDisease.list")
  -hit string
    	column name of hit to be filter (default "SampleID")
  -include string
    	include disease db file (default "/zfsyt1/B2C_RD_P2/USER/wangyaoshen/pipeline/filterExcel/etc/includeDisease.list")
  -input string
    	input excel
  -list string
    	hit list to be filter
  -output string
    	output excel, default -input.filter.xlsx
  -sep string
    	sep to split disease info (default "[n]")
  -sheet string
    	sheet name to be filter (default "All variants data")

```