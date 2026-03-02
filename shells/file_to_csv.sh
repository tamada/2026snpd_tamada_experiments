#! /bin/sh

shasum=$(sha256sum --quiet $1)
file=$(file -b $1)
size=$(stat -f %z $1)
ghidra_path=$(echo $1 | sed 's!build!ghidra!g')
ghidra_size=$(\ls -lR $ghidra_path | grep "^-" | awk '{sum += $5} END {print sum}')
items=$(echo $1 | tr '/' ',')

echo "$items,$shasum,\"$file\",$size,$ghidra_size"
