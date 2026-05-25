import pandas as pd
import glob
import os
import argparse
import sys
import numpy as np

def main():
    # Setup command line arguments
    parser = argparse.ArgumentParser(
        description="Extract pairs from similarity matrices in CSV files that exceed a specific threshold.",
        formatter_class=argparse.ArgumentDefaultsHelpFormatter
    )
    parser.add_argument("files", nargs="*", help="Path to target CSV files (wildcards supported)")
    parser.add_argument("-t", "--threshold", type=float, default=0.75, help="Threshold for similarity extraction")
    parser.add_argument("-o", "--output", help="Filename to save the results as a CSV")

    # Display help and exit if no arguments are provided
    if len(sys.argv) == 1 or not parser.parse_args().files:
        parser.print_help()
        sys.exit(0)

    args = parser.parse_args()
    results = []

    # Expand file paths (handles wildcards like *.csv)
    file_list = []
    for f_param in args.files:
        file_list.extend(glob.glob(f_param))

    if not file_list:
        print("Error: No target files found.", file=sys.stderr)
        sys.exit(1)

    for filepath in file_list:
        filename = os.path.basename(filepath)
        try:
            # Read CSV using the first column as row labels (index)
            df = pd.read_csv(filepath, index_col=0)
            
            # Convert the matrix to a stacked format (Row Label, Column Label, Value)
            stacked = df.stack()
            
            # Filter entries based on the threshold
            high_sim = stacked[stacked > args.threshold]
            
            # Iterate through the filtered results and store them
            for (row_label, col_label), similarity in high_sim.items():
                results.append({
                    "Row Label": row_label,
                    "Column Label": col_label,
                    "File Name": filename,
                    "Similarity": similarity
                })
        except Exception as e:
            print(f"Warning: Failed to process {filename}: {e}", file=sys.stderr)

    # Process final output
    if results:
        output_df = pd.DataFrame(results)
        if args.output:
            output_df.to_csv(args.output, index=False, encoding="utf-8-sig")
            print(f"Results successfully saved to {args.output}")
        else:
            # Print to standard output
            print(output_df.to_string(index=False))
    else:
        print(f"No data found exceeding the given threshold {args.threshold}.")

if __name__ == "__main__":
    main()