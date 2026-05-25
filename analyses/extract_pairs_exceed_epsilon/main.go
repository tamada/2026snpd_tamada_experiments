package main

import (
	"bufio"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	flag "github.com/spf13/pflag"
)

type Options struct {
	files     []string
	threshold float64
	dest      string
}

func helpMessage() string {
	return `Uasge: extract_pairs_exceed_epsilon [OPTIONS] <FILES...>
OPTIONS:
  -d, --dest <FILE>         Destination for output (default "-")
  -t, --threshold <VALUE>   Threshold for extracting pairs (default 0.75)
FILES
  Input files to process`
}

func ParseOptions(args []string) (*Options, error) {
	flags := flag.NewFlagSet("extract_pairs_exceed_epsilon", flag.ExitOnError)
	options := &Options{}
	flags.Usage = func() { fmt.Println(helpMessage()) }
	flags.StringVarP(&options.dest, "dest", "d", "-", "Destination for output")
	flags.Float64VarP(&options.threshold, "threshold", "t", 0.75, "Threshold for extracting pairs")
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

func processFile(file string, threshold float64, out *bufio.Writer) error {
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()
	csvio := csv.NewReader(f)
	csvio.FieldsPerRecord = -1 // Allow variable number of fields
	colHeaders, err := csvio.Read()
	if err != nil {
		return err
	}
	colHeaders = colHeaders[1:] // Skip the first column header
	errs := []error{}
	for {
		record, err := csvio.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			errs = append(errs, err)
			continue
		}
		rowHeader := record[0]
		for i, value := range record[1:] {
			value = strings.TrimSpace(value)
			if value == "" || strings.HasSuffix(value, "nsec)") {
				continue
			}
			value, err := strconv.ParseFloat(value, 64)
			if err != nil {
				errs = append(errs, err)
				continue
			}
			if value > threshold {
				fmt.Fprintf(out, "%s,%s,%s,%v,%f\n", rowHeader, colHeaders[i], file, rowHeader == colHeaders[i], value)
			}
		}
	}
	out.Flush()
	if len(errs) == 1 {
		return errs[0]
	} else if len(errs) > 1 {
		return errors.Join(errs...)
	} else {
		return nil
	}
}

func perform(options *Options) int {
	output, err := openDestination(options.dest)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening destination: %v\n", err)
		return 1
	}
	defer output.Close()
	errs := []error{}
	out := bufio.NewWriter(output)
	for _, file := range options.files {
		if err := processFile(file, options.threshold, out); err != nil {
			errs = append(errs, err)
		}
	}
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
