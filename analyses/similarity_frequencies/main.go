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
	files    []string
	interval float64
	dest     string
}

func (o *Options) intervals() []Interval {
	intervals := []Interval{}
	for lower := 0.0; lower < 1.0; lower += o.interval {
		upper := lower + o.interval
		if upper > 1.0 {
			upper = 1.0
		}
		name := fmt.Sprintf("[%.2f, %.2f)", lower, upper)
		if upper == 1.0 {
			name = fmt.Sprintf("[%.2f, %.2f]", lower, upper)
		}
		intervals = append(intervals, Interval{
			Lower: lower,
			Upper: upper,
			Name:  name,
		})
	}
	return intervals
}

func helpMessage() string {
	return `Uasge: similarity_frequencies [OPTIONS] <FILES...>
OPTIONS:
  -d, --dest <FILE>         Destination for output (default "-")
  -i, --interval <FLOAT>    Similarity interval of frequencies (default 0.1)
FILES
  Input files to process`
}

func ParseOptions(args []string) (*Options, error) {
	flags := flag.NewFlagSet("similarity_frequencies", flag.ExitOnError)
	options := &Options{}
	flags.Usage = func() { fmt.Println(helpMessage()) }
	flags.StringVarP(&options.dest, "dest", "d", "-", "Destination for output")
	flags.Float64VarP(&options.interval, "interval", "i", 0.1, "Similarity interval of frequencies")
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

type Interval struct {
	Lower float64
	Upper float64
	Name  string
}

func (i *Interval) Within(value float64) bool {
	return i.Lower <= value && value < i.Upper
}

func processFile(file string, intervals []Interval, out *bufio.Writer) error {
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()
	freq := map[string]int{}
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
			if rowHeader == colHeaders[i] {
				continue
			}
			value, err := strconv.ParseFloat(value, 64)
			if err != nil {
				errs = append(errs, err)
				continue
			}
			added := false
			for _, interval := range intervals {
				if interval.Within(value) {
					freq[interval.Name]++
					added = true
					break
				}
			}
			if !added {
				fmt.Printf("Value does not fit in any interval: %f, %s\n", value, file)
				freq[intervals[len(intervals)-1].Name]++
				added = true
			}
		}
	}
	fmt.Fprintf(out, "%s", file)
	for _, interval := range intervals {
		fmt.Fprintf(out, ",%d", freq[interval.Name])
	}
	fmt.Fprintln(out)
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
	intervals := options.intervals()
	for _, interval := range intervals {
		fmt.Fprintf(out, ",\"%s\"", interval.Name)
	}
	fmt.Fprintln(out)
	for _, file := range options.files {
		if err := processFile(file, intervals, out); err != nil {
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
