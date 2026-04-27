# parallel_coordinates_plot_generator
 
 Create a CSV file for parallel coordinates plots from the merged CSV file created by `merge_results`.
 The output CSV file has the following format:
 
```csv
Label,Language,Baseline,Diff_OS,Diff_Compiler,Diff_Arch
linux_go_amd64,go,1.0,0.847736,0.638777,0.915776
windows_rs-pure_arm64,rust,1.0,0.623582,0.984452,0.880337
linux_rs-sys_amd64,rust,1.0,0.808298,0.984705,0.794166
darwin_go_arm64,go,1.0,0.840673,0.590353,0.798273
darwin_rs-sys_arm64,rust,1.0,0.967443,0.984988,0.938754
linux_go_arm64,go,1.0,0.840673,0.649363,0.915776
```

## Results

In all of the parallel coordinates plots, the x-axis represents the different OS, compilers, and architectures.
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
