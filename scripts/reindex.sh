#!/bin/bash
# reindex操作由es服务器内部起多线程完成，效率比较高，90s可以完成2G数据的reindex.Works中450G的数据预估需要5.6个小时
st_time=$(date +%s)
# read -p "Source Index Name: " Source
# read -p "Destination Index Name: " Destination
# 读取密码,不显示密码
read -s -p "Password: " password
echo "Reindexing $Source to $Destination"
# 修改刷新间隔为-1，关闭自动刷新
echo "Disabling refresh interval"
curl -H "Content-Type: Application/json" -uelastic:$password -XPUT localhost:9200/$Destination/_settings -d '{
    "index": {
        "refresh_interval": "-1"
    }
}'
echo "Reindexing $Source to $Destination started"
# 设置reindex的线程数为auto，由es自动分配
curl -H "Content-Type: Application/json" -uelastic:$password -p -XPOST "localhost:9200/_reindex/?slices=auto&refresh"  -d '{
        "source": {
            "index": "works_v1",
            "size": 5000
        },
        "dest": {
            "index": "works_v2"
        }
    }'
echo "Reindexing $Source to $Destination completed"
end_time=$(date +%s)
echo "Total time taken: $((end_time-st_time)) seconds"
# 修改刷新间隔为1s
# echo "Enabling refresh interval"
# curl -H "Content-Type: Application/json" -uelastic:$password -XPUT localhost:9200/$Destination/_settings -d '{
#     "index": {
#         "refresh_interval": "30s"
#     }
# }'