#!/usr/bin/env sh

formats=$(
  for i in $(ls -d */); 
  do
    echo $i | sed 's/\///g'
  done
)

for f in $formats;
do
  schemer2 -in=$f:./$f/test -outputImage=test$f.png
done
