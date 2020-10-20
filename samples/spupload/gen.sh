#! /bin/bash

cd "${0%/*}"

rm -rf upload/
mkdir -p upload
cd upload

for n in {1..1000}; do
  progress=$(( n * 100 / 1000 ))
  echo -ne "> $progress %"\\r
  dd if=/dev/urandom of=file$( printf %03d "$n" ).bin bs=1 count=$(( RANDOM + 1024 )) > /dev/null 2>&1
done