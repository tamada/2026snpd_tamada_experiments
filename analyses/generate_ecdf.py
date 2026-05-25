import pandas as pd
import numpy as np
import glob
import os

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
    # Input directory
    input_dir = 'data/experiment1/hungarian'
    output_file = 'data/experiment1/ecdf_hungarian.csv'
    
    files = sorted(glob.glob(os.path.join(input_dir, '*.csv')))
    if not files:
        print(f"No CSV files found in {input_dir}")
        return

    data = {}
    for f in files:
        metric_name = os.path.basename(f).replace('.csv', '')
        print(f"Processing {metric_name}...")
        vals = extract_similarities(f)
        if vals.size > 0:
            data[metric_name] = vals
        else:
            print(f"Warning: No valid data found in {f}")

    if not data:
        print("No data collected.")
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
    
    # Save to CSV
    df_ecdf.to_csv(output_file, index=False)
    print(f"Successfully generated ECDF CSV: {output_file}")

if __name__ == "__main__":
    main()
