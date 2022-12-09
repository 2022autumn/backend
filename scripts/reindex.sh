#!/bin/bash

function reIndex {
    echo creating "$INDEX"_dummy from $INDEX
    curl -H "Content-Type: Application/json" -XPOST $ELHOST:$ELPORT/_reindex -d '{
        "source": {
            "index": '\"$INDEX\"'
        },
        "dest": {
            "index": '\"$INDEX'_dummy"
        }
    }'

    echo
    echo "$INDEX"_dummy create complete

    echo delete $INDEX
    curl -XDELETE $ELHOST:$ELPORT/$INDEX

    echo
    echo wait for indexing
    sleep 3
    echo moving back $INDEX from "$INDEX"_dummy
    curl -H "Content-Type: Application/json" -XPOST $ELHOST:$ELPORT/_reindex -d '{
        "source": {
            "index": '\"$INDEX'_dummy"
        },
        "dest": {
            "index": '\"$INDEX\"'
        }
    }'

    echo
    echo move back $INDEX complete

    echo cleaning up..
    echo removing "$INDEX"_dummy
    curl -XDELETE "$ELHOST:$ELPORT/$INDEX"_dummy

    echo
    echo all job done!
    curl $ELHOST:$ELPORT/_cat/indices
}

read -p "Elasticsearch IP: " ELHOST
read -p "Elasticsearch Port: " ELPORT
read -p "Index Name: " INDEX

read -p "Continue (y/n)?" choice
case "$choice" in
y | Y)
    echo "yes"
    reIndex
    ;;
n | N)
    echo "bye"
    exit 0
    ;;
*) echo "invalid" ;;
esac
