#!/bin/bash
# Create test data dir for bulk
mkdir -p /data/openalex/testdata/authors/
mkdir -p /data/openalex/testdata/institutions/
mkdir -p /data/openalex/testdata/works/
mkdir -p /data/openalex/testdata/concepts/
mkdir -p /data/openalex/testdata/venues/

# Copy test data from /data/openalex to /data/testdata, 5 files for each subdir
rm -f /data/openalex/testdata/authors/*
rm -f /data/openalex/testdata/institutions/*
rm -f /data/openalex/testdata/works/*
rm -f /data/openalex/testdata/concepts/*
rm -f /data/openalex/testdata/venues/*

# create real test data
echo "...Copying test data from /data/openalex to /data/openalex/testdata..."
echo "Copying authors..."
for i in {0..3}; do cp /data/openalex/authors/authors_data_$i.json /data/openalex/testdata/authors/;echo "finish copy authors_data_$i.json";done
echo "Copying institutions..."
for i in {0..5}; do cp /data/openalex/institutions/institutions_data_$i.json /data/openalex/testdata/institutions/;echo "finish copy institutions_data_$i.json"; done
echo "Copying works..."
for i in {0..5}; do cp /data/openalex/works/works_data_$i.json /data/openalex/testdata/works/;echo "finish copy works_data_$i.json"; done
echo "Copying concepts..."
for i in {0..5}; do cp /data/openalex/concepts/concepts_data_$i.json /data/openalex/testdata/concepts/;echo "finish copy concepts_data_$i.json"; done
echo "Copying venues..."
for i in {0..5}; do cp /data/openalex/venues/venues_data_$i.json /data/openalex/testdata/venues/;echo "finish copy venues_data_$i.json"; done

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
