# Experiment 1: Resemblance Evaluation

In the first experiment, we compare the birthmarks of the essentially near-identical software to confirm that they exhibit high similarity.
This evaluation is crucial for establishing a baseline before assessing the similarities of different software in subsequent experiments.
By "essentially identical," we refer to software that shares the same OS, architecture, programming language, and compiler, with only slight modifications.
For this purpose, we utilize binaries from different versions of bzip2 (versions 1.0.1 to 1.0.8), a widely used data compression utility.
These binaries were obtained from [the official website](https://sourceware.org/pub/bzip2/).
The differences of the source code between these versions are shown in Table 2, with [`cloc`](https://github.com/AlDanial/cloc).
Table 2 shows the differences between these versions of source code written in C,
where the "Versions" column indicates the two versions being compared.
And the "Type" column categorizes the differences into four types: same, modified, added, and removed,
for each of the following metrics: number of files, blank lines, comment lines, and code lines.
From table 2, we can see that the differences between these versions are quite small,
with only a few percent of code added, removed, or modified.

## Table 2. The differences between the versions of bzip2

| Versions | Type | Files | Blank line | Comments | Code |
| :--- | :--- | ---: | ---: | ---: | ---: |
| 1.0.1, 1.0.2 | Same | 2 | 26 | 1,104 | 4,879 |
| 1.0.1, 1.0.2 | Modified | 10 | 0 | 12 | 128 |
| 1.0.1, 1.0.2 | Added | 1 | 20 | 79 | 248 |
| 1.0.1, 1.0.2 | Removed | 0 | 0 | 10 | 115 |
| 1.0.2, 1.0.3 | Same | 4 | 46 | 1,175 | 5,224 |
| 1.0.2, 1.0.3 | Modified | 9 | 0 | 20 | 28 |
| 1.0.2, 1.0.3 | Added | 1 | 4 | 28 | 23 |
| 1.0.2, 1.0.3 | Removed | 0 | 0 | 0 | 3 |
| 1.0.3, 1.0.4 | Same | 0 | 0 | 761 | 5,226 |
| 1.0.3, 1.0.4 | Modified | 13 | 0 | 116 | 46 |
| 1.0.3, 1.0.4 | Added | 0 | 15 | 30 | 22 |
| 1.0.3, 1.0.4 | Removed | 0 | 75 | 346 | 3 |
| 1.0.4, 1.0.5 | Same | 1 | 17 | 883 | 5,290 |
| 1.0.4, 1.0.5 | Modified | 12 | 0 | 24 | 4 |
| 1.0.4, 1.0.5 | Added | 0 | 0 | 0 | 1 |
| 1.0.4, 1.0.5 | Removed | 0 | 0 | 0 | 0 |
| 1.0.5, 1.0.6 | Same | 1 | 17 | 883 | 5,292 |
| 1.0.5, 1.0.6 | Modified | 12 | 0 | 24 | 3 |
| 1.0.5, 1.0.6 | Added | 0 | 0 | 10 | 10 |
| 1.0.5, 1.0.6 | Removed | 0 | 0 | 0 | 0 |
| 1.0.6, 1.0.7 | Same | 1 | 17 | 893 | 5,285 |
| 1.0.6, 1.0.7 | Modified | 12 | 0 | 24 | 20 |
| 1.0.6, 1.0.7 | Added | 0 | 0 | 0 | 4 |
| 1.0.6, 1.0.7 | Removed | 0 | 0 | 0 | 0 |
| 1.0.7, 1.0.8 | Same | 1 | 17 | 893 | 5,295 |
| 1.0.7, 1.0.8 | Modified | 12 | 0 | 24 | 14 |
| 1.0.7, 1.0.8 | Added | 0 | 0 | 3 | 3 |
| 1.0.7, 1.0.8 | Removed | 0 | 0 | 0 | 0 |

## Directory layouts

- `box_plots.xlsx`: the data for plotting the box plots from the comparison results of hungarian, top-n-all, and top-n-one.
- `birthmarks`: extracted birthmarks of the binaries located on `executables` directory.
- `executables`: the binaries of bzip2 versions 1.0.1 to 1.0.8.
- `sources`: the source code of bzip2 versions 1.0.1 to 1.0.8 obtained from the official website.
- `hungarian`: the comparison results with Hungarian algorithm for aggregation of birthmarks.
  - `images`: the heatmap image from the tables directory data.
  - `results`: the raw comparison results in CSV format, listed the comparison pairs and their similarity scores.
  - `tables`: the comparison results in tabular format.
- `top-n-all`: the comparison results with top-n (n=all) algorithm for aggregation of birthmarks.
  - `images`: the heatmap image from the tables directory data.
  - `results`: the raw comparison results in CSV format, listed the comparison pairs and their similarity scores.
  - `tables`: the comparison results in tabular format.
- `top-n-one`: the comparison results with top-n (n=1) algorithm for aggregation of birthmarks.
  - `images`: the heatmap image from the tables directory data.
  - `results`: the raw comparison results in CSV format, listed the comparison pairs and their similarity scores.
  - `tables`: the comparison results in tabular format.
