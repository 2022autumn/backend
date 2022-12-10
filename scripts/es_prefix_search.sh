curl -XGET localhost:9200/_search -uelastic -p -H 'Content-Type: application/json' -d'
{
  "query": {
    "match_phrase_prefix": {
      "title": {
        "query": "Transformer"
      }
    }
  }
}'