package main

import (
	"fmt"
	"github.com/dustin/go-humanize"
)

func trackProgress(ids []string, progress <-chan Progress) {
	tracker := make(map[string]Progress)

	for p := range progress {
		tracker[p.Id] = p
		reportProgress(ids, tracker)
	}
	reportProgress(ids, tracker)
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
