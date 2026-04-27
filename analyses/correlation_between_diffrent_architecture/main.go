package main

import (
	"bufio"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	flag "github.com/spf13/pflag"
)

type Options struct {
	files []string
	dest  string
}

func helpMessage() string {
	return `Uasge: correlation_between_different_architectures [OPTIONS] <FILES...>
OPTIONS:
  -d, --dest <FILE>         Destination for output (default "-")
FILES
  Input files to process`
}

func ParseOptions(args []string) (*Options, error) {
	flags := flag.NewFlagSet("correlation_between_different_architectures", flag.ExitOnError)
	options := &Options{}
	flags.Usage = func() { fmt.Println(helpMessage()) }
	flags.StringVarP(&options.dest, "dest", "d", "-", "Destination for output")
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

type Target struct {
	Architecture    string
	OperatingSystem string
	Compiler        string
	Hash            string
}

func (t *Target) String() string {
	return fmt.Sprintf("%s_%s_%s", t.Architecture, t.OperatingSystem, t.Compiler)
}

func (t *Target) IsSameOSAndCompiler(other *Target) bool {
	return t.OperatingSystem == other.OperatingSystem && t.Compiler == other.Compiler
}

func ParseTarget(from string) (Target, error) {
	parts := strings.Split(from, "_")
	if len(parts) != 5 {
		return Target{}, fmt.Errorf("invalid target format: (%d): fields: %v", len(parts), from)
	}
	return Target{
		Architecture:    parts[1],
		OperatingSystem: parts[2],
		Compiler:        strings.TrimSuffix(parts[3], ".exe"),
		Hash:            strings.TrimSuffix(parts[4], ".json"),
	}, nil
}

type Comparison struct {
	Id         string
	Similarity float64
	TA         Target
	TB         Target
	Duration   float64
	Algorithm  string
}

func read_algorithm_from_filename(filename string) string {
	name := filepath.Base(filename)
	ext := filepath.Ext(name)
	return name[:len(name)-len(ext)-1]
}

func parseComparison(record []string, algorithm string) (Comparison, error) {
	if len(record) != 5 {
		return Comparison{}, fmt.Errorf("invalid record length: expected 5 fields (%d): %v", len(record), record)
	}
	similarity, err := strconv.ParseFloat(record[1], 64)
	if err != nil {
		return Comparison{}, fmt.Errorf("invalid similarity: %w", err)
	}
	duration, err := strconv.ParseFloat(record[4], 64)
	if err != nil {
		return Comparison{}, fmt.Errorf("invalid duration: %w", err)
	}
	ta, err := ParseTarget(record[2])
	if err != nil {
		return Comparison{}, fmt.Errorf("invalid target A: %w", err)
	}
	tb, err := ParseTarget(record[3])
	if err != nil {
		return Comparison{}, fmt.Errorf("invalid target B: %w", err)
	}
	// fmt.Fprintf(os.Stderr, "parsed comparison: id=%s, similarity=%f, ta=%v, tb=%v, duration=%f, algorithm=%s\n", record[0], similarity, ta, tb, duration, algorithm)

	return Comparison{
		Id:         record[0],
		Similarity: similarity,
		TA:         ta,
		TB:         tb,
		Duration:   duration,
		Algorithm:  algorithm,
	}, nil
}

func readComparison(file string) ([]Comparison, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	reader := csv.NewReader(f)
	reader.FieldsPerRecord = -1
	algorithm := read_algorithm_from_filename(file)
	reader.TrimLeadingSpace = true
	var comparisons []Comparison
	errs := []error{}
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			errs = append(errs, err)
			continue
		}
		if record[0] == "total duration" {
			continue
		}
		c, err := parseComparison(record, algorithm)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		comparisons = append(comparisons, c)
	}
	return comparisons, errors.Join(errs...)
}

func findOther(comparisons []Comparison, ta *Target, tb *Target, architecture string) (*Comparison, error) {
	for _, c := range comparisons {
		if c.TA.Architecture == c.TB.Architecture && c.TA.Architecture == architecture {
			if (c.TA.IsSameOSAndCompiler(ta) && c.TB.IsSameOSAndCompiler(tb)) ||
				(c.TA.IsSameOSAndCompiler(tb) && c.TB.IsSameOSAndCompiler(ta)) {
				return &c, nil
			}
		}
	}
	return nil, errors.New("other comparison not found")
}

func writeData(w io.Writer, comparisons []Comparison) []error {
	var errs []error
	fmt.Fprintf(os.Stderr, "found %d comparisons\n", len(comparisons))
	for _, c := range comparisons {
		if c.TA.Architecture == c.TB.Architecture && c.TA.Architecture == "arm64" {
			other, err := findOther(comparisons, &c.TA, &c.TB, "amd64")
			if err != nil {
				fmt.Fprintf(w, "%s,%s,,,%f,NOT_FOUND\n", c.TA.String(), c.TB.String(), c.Similarity)
			} else {
				fmt.Fprintf(w, "%s,%s,%s,%s,%f,%f,FOUND\n", c.TA.String(), c.TB.String(), other.TA.String(), other.TB.String(), c.Similarity, other.Similarity)
			}
		}
	}
	return errs
}

func perform(options *Options) int {
	output, err := openDestination(options.dest)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening destination: %v\n", err)
		return 1
	}
	defer output.Close()
	errs := []error{}
	var comparisons []Comparison
	for _, file := range options.files {
		c, err := readComparison(file)
		if err != nil {
			errs = append(errs, err)
			continue
		} else {
			comparisons = append(comparisons, c...)
		}
	}
	out := bufio.NewWriter(output)
	writeData(out, comparisons)
	out.Flush()

	if len(errs) > 0 {
		for _, err := range errs {
			fmt.Fprintf(os.Stderr, "Error processing file: %v\n", err)
		}
		return 1
	}
	return 0
}

func goMain(args []string) int {
	options, err := ParseOptions(args)
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
