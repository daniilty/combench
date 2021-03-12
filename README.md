# combench
***

## Purpose
Simple `go test -bench` results comparator

## Installation
```bash
$ go get -u github.com/daniilty/combench
```

## Usage
```bash
$ combench old.txt new.txt
```

## Example output
```bash
=> Difference in Total operations: new results(8627197) are differ from old (9088918) on -5.080044 %
   Difference in ns per operation: new results(147.9) are differ from old (127) on +16.456693 %
```