package main

import (
	"flag"
	"fmt"
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

	downloads := make([]Request, 0, flag.NArg()/2)

	for i := 0; i < flag.NArg(); i +=2 {
		downloads = append(downloads, Request{
			Url:      flag.Arg(i),
			Filename: flag.Arg(i+1),
		})
	}

	ids := make([]string, 0, len(downloads))


	for _, dl := range downloads {
		fmt.Println()
		ids = append(ids, dl.Filename)
	}

	fmt.Println("Starting download.")
	progress := StartBatch(*simDownloads, 100 * time.Millisecond, downloads...)
	trackProgress(ids, progress)
	fmt.Printf("Finished downloading in %s\n", time.Since(start))
}

