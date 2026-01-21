package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"os/exec"

	"github.com/artbin/mlist/mlist"
	"github.com/artbin/mlist/segment"
	"github.com/bradleyjkemp/memviz"
)

type listStruct[E any] struct {
	last    *segment.Segment[E]
	first   *segment.Segment[E]
	len     int
	horizon int
}

var outDir string

func init() {
	flag.StringVar(&outDir, "out", "visualize", "output directory")
	flag.Parse()

	err := os.MkdirAll(outDir, 0755)
	if err != nil {
		panic(err)
	}
}

func main() {
	fmt.Println("Visualizing MList (partially persistent)...")
	visualizeMList()

	fmt.Println("Done! Check the output directory:", outDir)
}

func visualizeMList() {
	list := mlist.MList[int]{}
	vizList := listStruct[int]{}

	// Empty list
	visualize("ilist", 0, vizList)

	// Push elements 1-8
	for value := 1; value <= 8; value++ {
		list = list.Push(value)
		exported := mlist.MListExport(list)
		vizList = listStruct[int]{
			last:    exported.Last,
			first:   exported.First,
			len:     exported.Len,
			horizon: exported.Horizon,
		}
		visualize("ilist", value, vizList)
	}

	// Pop elements
	for step := 1; step <= 8; step++ {
		list, _, _ = list.Pop()
		exported := mlist.MListExport(list)
		vizList = listStruct[int]{
			last:    exported.Last,
			first:   exported.First,
			len:     exported.Len,
			horizon: exported.Horizon,
		}
		visualize("ilist", 8+step, vizList)
	}
}

func visualize(prefix string, step int, list listStruct[int]) {
	buf := &bytes.Buffer{}

	memviz.Map(buf, &list)

	filename := fmt.Sprintf("%s/%s_%02d", outDir, prefix, step)
	err := os.WriteFile(filename+".dot", buf.Bytes(), 0644)
	if err != nil {
		panic(err)
	}

	png, err := exec.Command("dot", "-Tpng", filename+".dot").Output()
	if err != nil {
		fmt.Printf("Warning: failed to generate PNG for %s: %v\n", filename, err)
		return
	}

	err = os.WriteFile(filename+".png", png, 0644)
	if err != nil {
		panic(err)
	}
}
