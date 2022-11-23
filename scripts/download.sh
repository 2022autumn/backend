#!/bin/bash
# Download the data through the urls
# require: ./urls/*_url.txt preprocessed
# Usage: bash download.sh <target_dir>  e.g. bash ./download.sh ../data 1
# continue from the last downloaded file(if continue is 1)
# restart from the first file(if continue is 0)
if [ $# -lt 2 ]; then
    echo "Usage: bash download.sh <target_dir> <continue>"
    exit 1
fi
start=$(date +%s)
entity=('concepts' 'institutions' 'venues' 'works' 'authors')
target_dir=$1
mkdir -p $target_dir
if [ ! -d $target_dir ]; then
    echo "Target directory does not exist"
    exit 1
fi
# target_dir can not be a empty string 防止误操作删除根目录
if [ -z $target_dir ]; then
    echo "Target directory can not be a empty string"
    exit 1
fi
# if continue is 0, restart from the first file
if [ $2 -eq 0 ]; then
    echo "Restart from the first file"
    rm -rf $target_dir/*
    for e in ${entity[@]}; do
        mkdir -p $target_dir/$e
    done
fi
# if continue is 1, continue from the last downloaded file
if [ $2 -eq 1 ]; then
    echo "Continue from the last downloaded file"
    for e in ${entity[@]}; do
        mkdir -p $target_dir/$e
        # 前置处理 遍历$target_dir/$e
        for f in $(ls $target_dir/$e); do
            # 如果后缀为json，说明已经下载完成，跳过
            if [ ${f##*.} == "json" ]; then
                continue
            fi
            # 否则删除
            echo "Delete $f"
            # rm -rf $target_dir/$e/$f
            echo "Delete $target_dir/$e/$f"
        done
    done
fi
rm -rf ./error.log
for e in ${entity[@]}; do
    echo "-----Downloading $e..."
    number=0
    while read url; do
        # if countinue is 1, continue from the last downloaded file
        if [ $2 -eq 1 ]; then
            # if the file exists, continue
            if [ -f $target_dir/$e/${e}_data_${number}.json ]; then
                echo "-----$target_dir/$e/${e}_data_${number}.json exists, continue..."
                number=$((number + 1))
                continue
            fi
            # if the file filterred, continue
            if [ -f $target_dir/$e/filterred_${e}_data_${number}.json]; then
                echo "-----$target_dir/$e/filterred_${e}_data_${number}.json filterred, continue..."
                number=$((number + 1))
                continue
            fi
        fi
        start_time=$(date +%s)
        wget -P $target_dir/$e $url
        if [ $? -ne 0 ]; then
            echo "-----Download $url failed" >> error.log
            echo "-----Download $url failed"
            continue
        fi
        end_time=$(date +%s)
        basename=$(basename $url)
        nosuffix=${basename%.*}
        gzip -d $target_dir/$e/$basename
        # 把$nosuffix文件改名为${e}_data_${number}.json
        mv $target_dir/$e/$nosuffix $target_dir/$e/${e}_data_${number}.json
        echo "Downloaded $url in $((end_time - start_time)) seconds"
        echo "Made $target_dir/$e/${e}_data_${number}.json"
        number=$((number + 1))
    done <urls/${e}_url.txt
done
end=$(date +%s)
echo "-----Download finished."
runtime=$((end - start))
echo "runtime: $runtime seconds"
