import sys
import csv
import re

def parse_metadata(path):
    # birthmarks/experiment2/bzip2_amd64_darwin_clang_28dc2ee7.json
    filename = path.split('/')[-1]
    match = re.match(r"bzip2_([^_]+)_([^_]+)_([^_.\d]+)", filename)
    if not match:
        return None
    arch, os, compiler = match.groups()
    
    # 言語の判定
    lang = "C"
    if compiler in ["go", "tgo", "tinygo"]:
        lang = "Go"
    elif compiler in ["rs-pure", "rs-sys"]:
        lang = "Rust"
    
    return {"arch": arch, "os": os, "comp": compiler, "lang": lang}

def get_category(m1, m2):
    # if the languages are different, compilers are different, of course.
    if m1["lang"] != m2["lang"] and m1["os"] == m2["os"] and m1["arch"] == m2["arch"]:
        return "Cross-Language"
    
    if m1["lang"] == m2["lang"] and m1["comp"] != m2["comp"] and m1["os"] == m2["os"] and m1["arch"] == m2["arch"]:
        return "Cross-Compiler"
        
    if m1["lang"] == m2["lang"] and m1["comp"] == m2["comp"] and m1["os"] != m2["os"] and m1["arch"] == m2["arch"]:
        return "Cross-OS"
        
    if m1["lang"] == m2["lang"] and m1["comp"] == m2["comp"] and m1["os"] == m2["os"] and m1["arch"] != m2["arch"]:
        return "Cross-Arch"
#        if m1["os"] == "windows":
#            return "Cross-Arch (Windows)"
#        else:
#            return "Cross-Arch (Non-Windows)"

    return None

def main():
    if len(sys.argv) < 2:
        print("Usage: python generate_ecdf_csv.py <input_csv>", file=sys.stderr)
        sys.exit(1)

    categories = [
#        "Cross-Arch (Windows)",
#        "Cross-Arch (Non-Windows)",
        "Cross-Arch",
        "Cross-OS",
        "Cross-Compiler",
        "Cross-Language"
    ]
    data = {cat: [] for cat in categories}

    with open(sys.argv[1], 'r') as f:
        reader = csv.reader(f)
        for row in reader:
            if not row or len(row) < 4: continue
            try:
                sim = float(row[1])
                m1 = parse_metadata(row[2])
                m2 = parse_metadata(row[3])
                if not m1 or not m2: continue
                
                cat = get_category(m1, m2)
                if cat:
                    data[cat].append(sim)
            except ValueError:
                continue

    # グラフ描画用に 0.00 から 1.00 までの ECDF 値を算出
    x_steps = [i / 100.0 for i in range(101)]
    print("similarity," + ",".join(categories))
    
    for x in x_steps:
        row = [f"{x:.2f}"]
        for cat in categories:
            sims = data[cat]
            if not sims:
                row.append("0.000")
                continue
            count = sum(1 for s in sims if s <= x)
            row.append(f"{count / len(sims):.3f}")
        print(",".join(row))

if __name__ == "__main__":
    main()
