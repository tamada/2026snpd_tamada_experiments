import argparse
import pandas as pd
import numpy as np
import glob
import os
import sys

def extract_similarities(csv_path):
    # Read CSV, first column as index
    try:
        df = pd.read_csv(csv_path, index_col=0)
    except Exception as e:
        print(f"Error reading {csv_path}: {e}")
        return np.array([])

    # Convert to numeric, handle errors as NaN
    df = df.apply(pd.to_numeric, errors='coerce')
    
    # Drop rows/columns that are all NaN (like 'total duration' or empty headers)
    df = df.dropna(how='all', axis=0).dropna(how='all', axis=1)
    
    # Get values
    matrix = df.values
    
    # Check if square-ish (at least same number of rows and columns for the matrix part)
    # Actually, we just need the upper triangle.
    # If the matrix is not square, triu still works.
    
    # We want to extract pairs (i, j) where i < j.
    # This avoids diagonal (i=j) and duplicates (j, i).
    rows, cols = matrix.shape
    vals = []
    for i in range(rows):
        for j in range(i + 1, cols):
            val = matrix[i, j]
            if not np.isnan(val):
                vals.append(val)
    
    return np.array(vals)

def main():
    parser = argparse.ArgumentParser(
        description='Generate ECDF data from CSV files in a directory.'
    )
    parser.add_argument(
        'input_dir',
        nargs='?',
        help='Directory containing input CSV files (table format).'
    )
    parser.add_argument(
        '-o', '--output-file',
        help='Output CSV file. If omitted, write to standard output.'
    )
    args = parser.parse_args()

    if not args.input_dir:
        parser.print_help()
        return

    input_dir = args.input_dir
    output_file = args.output_file
    
    files = sorted(glob.glob(os.path.join(input_dir, '*.csv')))
    if not files:
        print(f"No CSV files found in {input_dir}", file=sys.stderr)
        return

    data = {}
    for f in files:
        metric_name = os.path.basename(f).replace('.csv', '')
        print(f"Processing {metric_name}...", file=sys.stderr)
        vals = extract_similarities(f)
        if vals.size > 0:
            data[metric_name] = vals
        else:
            print(f"Warning: No valid data found in {f}", file=sys.stderr)

    if not data:
        print("No data collected.", file=sys.stderr)
        return

    # Generate ECDF table
    # Similarity from 0.0 to 1.0
    similarity_levels = np.linspace(0, 1, 101) # 0.00, 0.01, ..., 1.00
    
    ecdf_results = {'Similarity': similarity_levels}
    
    for metric, vals in data.items():
        # ECDF calculation: for each x in similarity_levels, count how many vals <= x
        # Then divide by total number of vals
        n = len(vals)
        counts = [np.sum(vals <= x) / n for x in similarity_levels]
        ecdf_results[metric] = counts
    
    df_ecdf = pd.DataFrame(ecdf_results)
    
    # Save to CSV or standard output
    if output_file:
        df_ecdf.to_csv(output_file, index=False)
        print(f"Successfully generated ECDF CSV: {output_file}", file=sys.stderr)
    else:
        df_ecdf.to_csv(sys.stdout, index=False)

if __name__ == "__main__":
    main()
