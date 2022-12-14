# 别名处理
read -s -p "elastic Password: " password
curl -XGET "localhost:9200/_alias?pretty" -uelastic:$password -p -H 'Content-Type: application/json'
read -p "请输入要使用的新索引,如works_v2:" new_index
if [ -z $new_index ]; then
    echo "新索引不能为空"
    exit 1
fi
read -p "请输入alias指向的旧索引,如work_v1:,为空则表示不删除之前的alias" old_index
if [ -z $old_index ]; then
    echo "旧索引不能为空"
    curl -XPOST "localhost:9200/_aliases?pretty" -uelastic:$password -p -H 'Content-Type: application/json' -d'
    {
        "actions" : [
            { "add" : { "index" : "'$new_index'", "alias" : "works" } }
        ]
    }'
    curl -XGET "localhost:9200/_alias?pretty" -uelastic:$password -p -H 'Content-Type: application/json'
    exit
fi
curl -XPOST "localhost:9200/_aliases?pretty" -uelastic:$password -p -H 'Content-Type: application/json' -d'
{
    "actions" : [
        { "remove" : { "index" : "'$old_index'", "alias" : "works" } },
        { "add" : { "index" : "'$new_index'", "alias" : "works" } }
    ]
}'
curl -XGET "localhost:9200/_alias?pretty" -uelastic:$password -p -H 'Content-Type: application/json'
