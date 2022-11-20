.PHONY: filter clean
filter:
	cat /data/openalex/works/works_data_20.json > ./work.json
	cat /data/openalex/authors/authors_data_20.json > ./author.json
	cat /data/openalex/institutions/institutions_data_1.json > ./institution.json
	cat /data/openalex/concepts/concepts_data_1.json > ./concept.json
	cat /data/openalex/venues/venues_data_1.json > ./venue.json

clean:
	rm -f *.json