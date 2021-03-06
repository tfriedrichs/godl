package main

import (
	"flag"
	"fmt"
	"github.com/tfriedrichs/godl"
	"os"
	"time"
)

func main() {
	start := time.Now()

	simDownloads := flag.Int("n", 3, "Number of simultaneous downloads")
	flag.Parse()

	if (flag.NArg() % 2) != 0 {
		fmt.Println("Illegal number of arguments.")
		os.Exit(1)
	}

	downloads := make([]godl.Request, 0, flag.NArg()/2)

	for i := 0; i < flag.NArg(); i +=2 {
		downloads = append(downloads, godl.Request{
			Url:      flag.Arg(i),
			Filename: flag.Arg(i+1),
		})
	}

	ids := make([]string, 0, len(downloads))

	fmt.Println("Starting download.")
	for _, dl := range downloads {
		fmt.Println()
		ids = append(ids, dl.Filename)
	}

	progress, err := godl.StartBatch(*simDownloads, 100 * time.Millisecond, downloads...)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	godl.TrackProgress(ids, progress)
	fmt.Printf("Finished downloading in %s\n", time.Since(start))
}
