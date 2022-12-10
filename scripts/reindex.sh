#!/bin/bash
# reindex操作由es服务器内部起多线程完成，效率比较高，90s可以完成2G数据的reindex.Works中450G的数据预估需要5.6个小时
st_time=$(date +%s)
read -p "Source Index Name: " Source
read -p "Destination Index Name: " Destination
echo "Reindexing $Source to $Destination"
curl -H "Content-Type: Application/json" -XPOST localhost:9200/_reindex -uelastic -p -d '{
        "source": {
            "index": '\"$Source\"'
        },
        "dest": {
            "index": '\"$Destination\"'
        }
    }'
echo "Reindexing $Source to $Destination completed"
end_time=$(date +%s)
echo "Total time taken: $((end_time-st_time)) seconds"