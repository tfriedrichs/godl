package main

import (
	"flag"
	"fmt"
	"github.com/dustin/go-humanize"
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

	fmt.Println("Starting download.")
	done, progress := StartBatch(*simDownloads, 100 * time.Millisecond, downloads...)
	ids := make([]string, 0, len(downloads))


	for _, dl := range downloads {
		fmt.Println()
		ids = append(ids, dl.Filename)
	}
	trackProgress(ids, progress, done)

	fmt.Printf("Finished downloading in %s\n", time.Since(start))
}

func trackProgress(ids []string, progress <-chan Progress, done <-chan struct{}) {
	tracker := make(map[string]Progress)
	for {
		select {
		case p := <- progress:
			tracker[p.Id] = p
			reportProgress(ids, tracker)
		case <-done:
			reportProgress(ids, tracker)
			return
		}
	}
}

func reportProgress(ids []string, dls map[string]Progress) {
	fmt.Printf("\u001b[%dA", len(ids))
	for _, id := range ids {
		p := dls[id]
		if p.Error != nil {
			fmt.Printf("\u001b[2K%s: Error (%s)\n", p.Id, p.Error)
		} else if p.Current == 0 {
			fmt.Printf("\u001B[2K%s: Waiting to start...\n", id)
		} else if p.Current == p.Total {
			fmt.Printf("\u001B[2K%s: Finished downloading %s in %s.\n", p.Id, humanize.Bytes(uint64(p.Total)), p.Elapsed)
		} else {
			fmt.Printf("\u001B[2K%s: Downloading (%s/%s).\n", p.Id, humanize.Bytes(uint64(p.Current)), humanize.Bytes(uint64(p.Total)))
		}
	}
}

