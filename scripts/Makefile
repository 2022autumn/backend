.PHONY: filter clean indexs bulk

# 创建测试文件
filter:
	bash create_testdata.sh

# 创建index
indexs:
	bash create_indexs.sh

reindex:
	bash reindex.sh

# 执行bulker
bulk:
	go run bulker.go > bulk.log 2>bulkproc.log

clean:
	rm -f *.json