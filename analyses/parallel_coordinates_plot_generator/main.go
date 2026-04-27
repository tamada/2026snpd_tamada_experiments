package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"

	flag "github.com/spf13/pflag"
)

type Options struct {
	files     []string
	algorithm string
	dest      string
}

type Target struct {
	Architecture    string
	OperatingSystem string
	Compiler        string
	Language        string
}

func NewTarget(arch, os, compiler string) *Target {
	language := "unknown"
	c := strings.ToLower(compiler)
	if strings.Contains(c, "gcc") || strings.Contains(c, "clang") || strings.Contains(c, "msvc") {
		language = "c"
	} else if strings.Contains(c, "go") || strings.Contains(c, "tgo") {
		language = "go"
	} else if strings.Contains(c, "rs-") {
		language = "rust"
	}
	return &Target{
		Architecture:    arch,
		OperatingSystem: os,
		Compiler:        compiler,
		Language:        language,
	}
}

func (t *Target) String() string {
	return fmt.Sprintf("%s_%s_%s", t.OperatingSystem, t.Compiler, t.Architecture)
}

func (t *Target) IsSame(other *Target) bool {
	return t.Architecture == other.Architecture &&
		t.OperatingSystem == other.OperatingSystem &&
		t.Compiler == other.Compiler
}

type Comparison struct {
	TA           *Target
	TB           *Target
	Similarities []float64
}

// 評価軸の判定関数群
func isDiffOs(a, b *Target) bool {
	return a.OperatingSystem != b.OperatingSystem && a.Compiler == b.Compiler && a.Architecture == b.Architecture
}
func isDiffComp(a, b *Target) bool {
	return a.OperatingSystem == b.OperatingSystem && a.Compiler != b.Compiler && a.Architecture == b.Architecture
}
func isDiffArch(a, b *Target) bool {
	return a.OperatingSystem == b.OperatingSystem && a.Compiler == b.Compiler && a.Architecture != b.Architecture
}

func main() {
	options, err := parseOptions(os.Args[1:])
	if err != nil || len(options.files) == 0 {
		fmt.Println(`Usage: parallel_coordinate_builder [OPTIONS] <FILE>
OPTIONS:
  -d, --dest <FILE>         Destination for output (default "-")
  -a, --algorithm <ALGO>    Algorithm to use (default: cosine)
FILE
  Input files to process (merged csv files from merge_results)`)
		return
	}

	data, _, err := readAll(options.files[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	algoIdx, _ := algorithmToIndex(options.algorithm)

	// 全バイナリ（Target）をリストアップ
	targetMap := make(map[string]*Target)
	for _, c := range data {
		targetMap[c.TA.String()] = c.TA
		targetMap[c.TB.String()] = c.TB
	}

	out, _ := openDestination(options.dest)
	defer out.Close()
	writer := bufio.NewWriter(out)
	defer writer.Flush()

	// Excel用ヘッダー
	// Baselineは常に1.0 (自分自身との比較)
	fmt.Fprintln(writer, "Label,Language,Baseline,Diff_OS,Diff_Compiler,Diff_Arch")

	for _, baseTarget := range targetMap {
		// 各軸の最大類似度を保持（同じ条件のペアが複数ある場合、最も良い結果を採用）
		maxOs, maxComp, maxArch := 0.0, 0.0, 0.0

		for _, c := range data {
			var other *Target
			if c.TA.IsSame(baseTarget) {
				other = c.TB
			} else if c.TB.IsSame(baseTarget) {
				other = c.TA
			} else {
				continue
			}

			sim := c.Similarities[algoIdx]
			if isDiffOs(baseTarget, other) && sim > maxOs {
				maxOs = sim
			} else if isDiffComp(baseTarget, other) && sim > maxComp {
				maxComp = sim
			} else if isDiffArch(baseTarget, other) && sim > maxArch {
				maxArch = sim
			}
		}

		// 全ての軸にデータが存在する場合のみ出力（グラフの線が途切れないようにするため）
		if maxOs > 0 && maxComp > 0 && maxArch > 0 {
			fmt.Fprintf(writer, "%s,%s,1.0,%f,%f,%f\n",
				baseTarget.String(), baseTarget.Language, maxOs, maxComp, maxArch)
		}
	}
}

// --- 以下、既存のヘルパー関数群 (一部流用) ---

func parseOptions(args []string) (*Options, error) {
	flags := flag.NewFlagSet("parallel_builder", flag.ExitOnError)
	options := &Options{}
	flags.StringVarP(&options.dest, "dest", "d", "-", "Output CSV")
	flags.StringVarP(&options.algorithm, "algorithm", "a", "simpson", "Algorithm")
	flags.Parse(args)
	options.files = flags.Args()
	return options, nil
}

func readAll(file string) ([]*Comparison, []string, error) {
	f, _ := os.Open(file)
	defer f.Close()
	r := csv.NewReader(f)
	header, _ := r.Read()
	var results []*Comparison
	for {
		rec, err := r.Read()
		if err == io.EOF {
			break
		}
		// 期待される列数: ID(0), Arch1(1), OS1(2), Comp1(3), Arch2(4), OS2(5), Comp2(6), Sim...
		if len(rec) < 8 {
			continue
		}
		ta := NewTarget(rec[1], rec[2], rec[3])
		tb := NewTarget(rec[4], rec[5], rec[6])
		sims := make([]float64, len(header)-7)
		for i := 7; i < len(rec); i++ {
			fmt.Sscanf(rec[i], "%f", &sims[i-7])
		}
		results = append(results, &Comparison{TA: ta, TB: tb, Similarities: sims})
	}
	return results, header, nil
}

func algorithmToIndex(algo string) (int, error) {
	indices := map[string]int{"cosine": 0, "dice": 1, "euclidean": 2, "jaccard": 3, "lcs": 4, "levenshtein": 5, "simpson": 6, "weighted_jaccard": 7}
	if idx, ok := indices[strings.ToLower(algo)]; ok {
		return idx, nil
	}
	return 0, nil
}

func openDestination(dest string) (io.WriteCloser, error) {
	if dest == "-" {
		return os.Stdout, nil
	}
	return os.Create(dest)
}
