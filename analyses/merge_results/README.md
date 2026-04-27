# merge_results

Merge the oinkie comparison result files into a single CSV file.
The output CSV file has the following format:

```csv
ID,similarity,A,B,duration,algorithm,identical,A,,,,B,,,
0,1,amd64_darwin_clang_28dc2ee7.json,amd64_darwin_gcc_885ac7de.json,1192998750,cosine,FALSE,amd64,darwin,clang,28dc2ee7.json,amd64,darwin,gcc,885ac7de.json
1,0.024638437,amd64_darwin_clang_28dc2ee7.json,amd64_darwin_go_cd368ba0.json,22177455416,cosine,FALSE,amd64,darwin,clang,28dc2ee7.json,amd64,darwin,go,cd368ba0.json
2,0.017551594077312233,amd64_darwin_clang_28dc2ee7.json,amd64_darwin_rs-pure_672b9aab.json,35848995125,cosine,FALSE,amd64,darwin,clang,28dc2ee7.json,amd64,darwin,rs-pure,672b9aab.json
3,0.017794800123211024,amd64_darwin_clang_28dc2ee7.json,amd64_darwin_rs-sys_089b3292.json,35827736666,cosine,FALSE,amd64,darwin,clang,28dc2ee7.json,amd64,darwin,rs-sys,089b3292.json
4,0.13012927159292503,amd64_darwin_clang_28dc2ee7.json,amd64_darwin_tgo_e9ce53c1.json,4732307708,cosine,FALSE,amd64,darwin,clang,28dc2ee7.json,amd64,darwin,tgo,e9ce53c1.json
```

The first line is the header.
Each line after the first line represents a comparison between two targets.

## Usage

