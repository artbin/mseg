package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

type row struct {
	pkg       string
	bench     string
	iters     string
	nsPerOp   string
	bytesPer  string
	allocsPer string
	mbPerS    string
}

var (
	benchLineRe = regexp.MustCompile(`^(Benchmark[^\s]+)\s+(\d+)\s+([0-9.]+)\s+ns/op(?:\s+([0-9.]+)\s+B/op\s+([0-9.]+)\s+allocs/op)?(?:\s+([0-9.]+)\s+MB/s)?$`)
)

func main() {
	var (
		count     = flag.Int("count", 5, "benchmark repetition count (go test -count)")
		benchtime = flag.String("benchtime", "100ms", "per-benchmark run time (go test -benchtime)")
		mods      = flag.String("mods", "mlist,pflist", "comma-separated module dirs to benchmark (relative to -pkg)")
		out       = flag.String("out", "", "output CSV path")
		pkg       = flag.String("pkg", ".", "repo root path (expects module dirs like ./mlist, ./pflist, ...)")
	)
	flag.Parse()

	if *out == "" {
		fmt.Fprintln(os.Stderr, "missing required -out")
		os.Exit(2)
	}

	repoRoot, err := filepath.Abs(*pkg)
	if err != nil {
		fmt.Fprintln(os.Stderr, "bad -pkg:", err)
		os.Exit(2)
	}

	modDirs := []string{}
	for _, m := range strings.Split(*mods, ",") {
		m = strings.TrimSpace(m)
		if m == "" {
			continue
		}
		modDirs = append(modDirs, m)
	}
	var rows []row

	for _, d := range modDirs {
		modPath := filepath.Join(repoRoot, d)
		if _, err := os.Stat(filepath.Join(modPath, "go.mod")); err != nil {
			continue
		}

		cmd := exec.Command("go", "test", "./...",
			"-run", "^$",
			"-bench", ".",
			"-benchmem",
			"-benchtime", *benchtime,
			"-count", strconv.Itoa(*count),
		)
		cmd.Dir = modPath
		// Keep stderr visible in case benchmarks fail to compile.
		outBytes, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Fprintf(os.Stderr, "bench failed in %s:\n%s\n", d, string(outBytes))
			os.Exit(1)
		}

		rows = append(rows, parseBenchOutput(d, outBytes)...)
	}

	if err := writeCSV(*out, rows); err != nil {
		fmt.Fprintln(os.Stderr, "write CSV:", err)
		os.Exit(1)
	}
}

func parseBenchOutput(pkg string, out []byte) []row {
	var rows []row
	s := bufio.NewScanner(strings.NewReader(string(out)))
	for s.Scan() {
		line := strings.TrimSpace(s.Text())
		if !strings.HasPrefix(line, "Benchmark") {
			continue
		}

		m := benchLineRe.FindStringSubmatch(line)
		if m == nil {
			// Ignore unknown benchmark line formats (e.g. benchmark status lines).
			continue
		}

		r := row{
			pkg:     pkg,
			bench:   m[1],
			iters:   m[2],
			nsPerOp: m[3],
		}
		if m[4] != "" {
			r.bytesPer = m[4]
			r.allocsPer = m[5]
		} else {
			r.bytesPer = ""
			r.allocsPer = ""
		}
		r.mbPerS = m[6]

		rows = append(rows, r)
	}
	return rows
}

func writeCSV(path string, rows []row) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	w := csv.NewWriter(f)
	defer w.Flush()

	// Keep the paper-friendly columns, but also include pkg+bench for easier plotting.
	if err := w.Write([]string{
		"name",
		"pkg",
		"bench",
		"iters",
		"ns_per_op",
		"bytes_per_op",
		"allocs_per_op",
		"mb_per_s",
	}); err != nil {
		return err
	}

	for _, r := range rows {
		name := r.pkg + "/" + r.bench
		if err := w.Write([]string{
			name,
			r.pkg,
			r.bench,
			r.iters,
			r.nsPerOp,
			r.bytesPer,
			r.allocsPer,
			r.mbPerS,
		}); err != nil {
			return err
		}
	}
	return w.Error()
}

