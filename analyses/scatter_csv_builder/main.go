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
	switch compiler {
	case "gcc", "clang", "msvc":
		language = "c"
	case "go", "tgo":
		language = "go"
	case "rs-pure", "rs-sys":
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
	return fmt.Sprintf("%s_%s_%s", t.Architecture, t.OperatingSystem, t.Compiler)
}

func (t *Target) IsSame(other *Target) bool {
	return t.Architecture == other.Architecture &&
		t.OperatingSystem == other.OperatingSystem &&
		t.Compiler == other.Compiler
}

func (t *Target) IsOnlyDifferentCompiler(other *Target) bool {
	return t.Architecture == other.Architecture &&
		t.OperatingSystem == other.OperatingSystem &&
		t.Language == other.Language &&
		t.Compiler != other.Compiler
}

func (t *Target) IsOnlyDifferentOs(other *Target) bool {
	return t.Architecture == other.Architecture &&
		t.OperatingSystem != other.OperatingSystem &&
		t.Language == other.Language &&
		t.Compiler == other.Compiler
}

type Comparison struct {
	ID           string
	TA           *Target
	TB           *Target
	Similarities []float64
}

func (c *Comparison) IsXItem() bool {
	return c.TA.IsOnlyDifferentOs(c.TB)
}

func algorithmToIndex(algorithm string) (int, error) {
	switch strings.ToLower(algorithm) {
	case "cosine":
		return COSINE_INDEX, nil
	case "dice":
		return DICE_INDEX, nil
	case "euclidean":
		return EUCLIDEAN_INDEX, nil
	case "jaccard":
		return JACCARD_INDEX, nil
	case "lcs":
		return LCS_INDEX, nil
	case "levenshtein":
		return LEVENSHTEIN_INDEX, nil
	case "simpson":
		return SIMPSON_INDEX, nil
	case "weighted_jaccard":
		return WEITHTED_JACCARD_INDEX, nil
	default:
		return -1, fmt.Errorf("%s: unknown algorithm", algorithm)
	}
}

func findOther(target *Target, comparison *Comparison) *Target {
	if comparison.TA.IsSame(target) {
		return comparison.TB
	} else if comparison.TB.IsSame(target) {
		return comparison.TA
	}
	return nil
}

func printLine(target *Target, comparison *Comparison, vsComparison *Comparison, index int, out io.Writer) {
	other := findOther(target, vsComparison)
	x := comparison.Similarities[index]
	y := vsComparison.Similarities[index]
	switch other.Language {
	case "c":
		fmt.Fprintf(out, "%f,%f,,,%s\n", x, y, target.String())
	case "go":
		fmt.Fprintf(out, "%f,,%f,,%s\n", x, y, target.String())
	case "rust":
		fmt.Fprintf(out, "%f,,,%f,%s\n", x, y, target.String())
	}
}

func findOnlyDifferentCompilerItem(data []*Comparison, target *Target) []*Comparison {
	var result []*Comparison
	for _, c := range data {
		if c.TA.IsSame(target) {
			if c.TB.IsOnlyDifferentCompiler(target) {
				result = append(result, c)
			}
		} else if c.TB.IsSame(target) {
			if c.TA.IsOnlyDifferentCompiler(target) {
				result = append(result, c)
			}
		}
	}
	return result
}

func processData(data []*Comparison, algorithm string, out io.Writer) error {
	index, err := algorithmToIndex(algorithm)
	if err != nil {
		return err
	}
	for _, comparison := range data {
		if comparison.IsXItem() {
			vsTaComparison := findOnlyDifferentCompilerItem(data, comparison.TA)
			if len(vsTaComparison) > 0 {
				for _, vs := range vsTaComparison {
					printLine(comparison.TA, comparison, vs, index, out)
				}
			}
			vsTbComparison := findOnlyDifferentCompilerItem(data, comparison.TB)
			if len(vsTbComparison) > 0 {
				for _, vs := range vsTbComparison {
					printLine(comparison.TB, comparison, vs, index, out)
				}
			}
		}
	}
	return nil
}

const (
	COSINE_INDEX           = 0
	DICE_INDEX             = 1
	EUCLIDEAN_INDEX        = 2
	JACCARD_INDEX          = 3
	LCS_INDEX              = 4
	LEVENSHTEIN_INDEX      = 5
	SIMPSON_INDEX          = 6
	WEITHTED_JACCARD_INDEX = 7
)

func helpMessage() string {
	return `Uasge: scatter_csv_builder [OPTIONS] <FILE>
OPTIONS:
  -d, --dest <FILE>         Destination for output (default "-")
  -a, --algorithm <ALGO>    Algorithm to use (default: cosine)
FILE
  Input files to process (merged csv files from merge_results)`
}

func parseOptions(args []string) (*Options, error) {
	flags := flag.NewFlagSet("scatter_csv_builder", flag.ExitOnError)
	options := &Options{}
	flags.Usage = func() { fmt.Println(helpMessage()) }
	flags.StringVarP(&options.dest, "dest", "d", "-", "Destination for output")
	flags.StringVarP(&options.algorithm, "algorithm", "a", "cosine", "Algorithm to use (default: cosine)")
	if err := flags.Parse(args); err != nil {
		return nil, err
	}
	options.files = flags.Args()
	return options, nil
}

func openDestination(dest string) (io.WriteCloser, error) {
	if dest == "-" {
		return os.Stdout, nil
	} else {
		file, err := os.Create(dest)
		if err != nil {
			return nil, err
		}
		return file, nil
	}
}

func parseLine(record []string, header []string) (Comparison, error) {
	id := record[0]
	ta := NewTarget(record[1], record[2], record[3])
	tb := NewTarget(record[4], record[5], record[6])
	similarities := make([]float64, len(header)-7)
	for i := 7; i < len(record); i++ {
		var sim float64
		_, err := fmt.Sscanf(record[i], "%f", &sim)
		if err != nil {
			return Comparison{}, fmt.Errorf("error parsing similarity value: %v", err)
		}
		similarities[i-7] = sim
	}
	return Comparison{ID: id, TA: ta, TB: tb, Similarities: similarities}, nil
}

func readAll(file string) ([]*Comparison, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	csvio := csv.NewReader(f)
	csvio.FieldsPerRecord = -1 // Allow variable number of fields
	results := []*Comparison{}
	header, err := csvio.Read()
	if err != nil {
		return nil, fmt.Errorf("error reading header: %v", err)
	}
	for {
		record, err := csvio.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("error reading record: %v", err)
		}
		comparison, err := parseLine(record, header)
		if err != nil {
			return nil, fmt.Errorf("error parsing line: %v", err)
		}
		results = append(results, &comparison)
	}

	return results, nil
}

func perform(options *Options) int {
	output, err := openDestination(options.dest)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening destination: %v\n", err)
		return 1
	}
	defer output.Close()
	data, err := readAll(options.files[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		return 1
	}
	out := bufio.NewWriter(output)
	defer out.Flush()
	err = processData(data, options.algorithm, out)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error processing data: %v\n", err)
		return 1
	}
	return 0
}

func goMain(args []string) int {
	options, err := parseOptions(args)
	if err != nil {
		return 1
	}
	if len(options.files) == 0 {
		return 1
	}
	return perform(options)
}

func main() {
	status := goMain(os.Args[1:])
	os.Exit(status)
}
