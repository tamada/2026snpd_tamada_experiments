# scatter_csv_builder

Create a CSV file for scatter plots from the merged CSV file created by `merge_results`.
The output CSV file has the following format:

```csv
<ALGORITHM>,C,Go,Rust,base_binary
0.0.98,,,0.81,bzip2_amd64_darwin_rs-pure_672b9aab.json
```

## Results

In all of the following scatter plots, the x-axis represents the similarity score between the different OS.
The y-axis represents the similarity score between the different compilers/features.

### cosine

![cosine](images/cosine.pdf)

### dice

![dice](images/dice.pdf)

### euclidean

![euclidean](images/euclidean.pdf)

### jaccard

![jaccard](images/jaccard.pdf)

### lcs

![lcs](images/lcs.pdf)

### levenshtein

![levenshtein](images/levenshtein.pdf)

### simpson

![simpson](images/simpson.pdf)

### weighted-jaccard

![weighted-jaccard](images/weighted-jaccard.pdf)
