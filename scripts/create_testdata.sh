#!/bin/bash
# Create test data dir for bulk
mkdir -p /data/testdata/authors/
mkdir -p /data/testdata/institutions/
mkdir -p /data/testdata/works/
mkdir -p /data/testdata/concepts/
mkdir -p /data/testdata/venues/

# Copy test data from /data/openalex to /data/testdata, 5 files for each subdir
rm -f /data/testdata/authors/*
rm -f /data/testdata/institutions/*
rm -f /data/testdata/works/*
rm -f /data/testdata/concepts/*
rm -f /data/testdata/venues/*

# create real test data
echo "...Copying test data from /data/openalex to /data/testdata..."
echo "Copying authors..."
for i in {0..4}; do cp /data/openalex/authors/filterred_authors_data_$i.json /data/testdata/authors/;echo "finish copy filterred_authors_data_$i.json";done
echo "Copying institutions..."
for i in {0..4}; do cp /data/openalex/institutions/filterred_institutions_data_$i.json /data/testdata/institutions/;echo "finish copy filterred_institutions_data_$i.json"; done
echo "Copying works..."
for i in {0..4}; do cp /data/openalex/works/filterred_works_data_$i.json /data/testdata/works/;echo "finish copy filterred_works_data_$i.json"; done
echo "Copying concepts..."
for i in {0..4}; do cp /data/openalex/concepts/filterred_concepts_data_$i.json /data/testdata/concepts/;echo "finish copy filterred_concepts_data_$i.json"; done
echo "Copying venues..."
for i in {0..4}; do cp /data/openalex/venues/filterred_venues_data_$i.json /data/testdata/venues/;echo "finish copy filterred_venues_data_$i.json"; done

# craete null test file
# echo "...Copying test data from /data/openalex to /data/testdata..."
# echo "Copying authors..."
# for i in {0..4}; do touch /data/testdata/authors/filterred_authors_data_$i.json;done
# echo "Copying institutions..."
# for i in {0..4}; do touch /data/testdata/institutions/filterred_institutions_data_$i.json; done
# echo "Copying works..."
# for i in {0..4}; do touch /data/testdata/works/filterred_works_data_$i.json; done
# echo "Copying concepts..."
# for i in {0..4}; do touch /data/testdata/concepts/filterred_concepts_data_$i.json; done
# echo "Copying venues..."
# for i in {0..4}; do touch /data/testdata/venues/filterred_venues_data_$i.json; done
