#!/bin/bash
# Download the data through the urls
# require: ./urls/*_url.txt preprocessed
# Usage: bash download.sh <target_dir>  e.g. bash ./download.sh ../data
if [ $# -lt 1 ]; then
    echo "Usage: bash download.sh <target_dir>"
    exit 1
fi
entity=('authors' 'concepts' 'works' 'institutions' 'venues')
target_dir=$1
mkdir -p $target_dir
rm -rf $target_dir
for e in ${entity[@]}; do
    echo "-----Downloading $e..."
    mkdir -p $target_dir/$e
    rm -f $target_dir/$e/*
    while read url; do
        echo "Downloading $url"
        wget -P $target_dir/$e $url
        basename=$(basename $url)
        nosuffix=${basename%.*}
        gzip -d $target_dir/$e/$basename
        cat $target_dir/$e/$nosuffix >>$target_dir/$e/${e}_data.json
        rm -f $target_dir/$e/$nosuffix
    done <urls/${e}_url.txt
done
echo "-----Download finished."
