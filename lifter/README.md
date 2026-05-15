# Binary Lifter with Ghidra

This directory contains a script `lift.sh` that performs binary lifting using Ghidra. The script takes a target binary as an argument and generates a JSON file containing the lifted PCode representation of the binary.

## :runner: Usage

To use the `lift.sh` script, run the following command in your terminal:

```sh
lifter/lift.sh /path/to/target/binary
```

The `lift.sh` assumes that:

- the working directory is the root of the project (the parent directory of `lifter`), and
- the path of binary files contains `executables`.

The script is to run Ghidra with headless mode (without GUI).

## the Generated JSON

### Examples

```json
{
    "program": "factorizer", // The name of the lifted program
    "path": "/home/tamada/"
}
```
