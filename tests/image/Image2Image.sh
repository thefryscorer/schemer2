#!/usr/bin/env sh

if (( $# < 2 ))  
then
  echo "Error, missing parameters. Example: ./Image2Image.sh width height"
else
  for i in `seq 1 10`;
  do
    schemer2 -minBright=0 -in=img:testinput.png -outputImage=testout$i.png -w=$1 -h=$2
  done
fi
