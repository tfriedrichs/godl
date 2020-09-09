package main

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type Request struct {
	Url string
	Filename string
}

func StartBatch(simDownloads int, reportInterval time.Duration, requests ...Request) <-chan Progress {

	wg := &sync.WaitGroup{}
	work := make(chan Request, len(requests))
	progress := make(chan Progress, 3)
	wg.Add(len(requests))

	for _, req := range requests {
		work <- req
	}

	go func() {
		defer close(progress)
		wg.Wait()
	}()

	for i := 0; i < simDownloads; i++ {
		go workDownload(work, progress, reportInterval, wg)
	}

	return progress
}

func workDownload(requests <-chan Request, progress chan<- Progress, interval time.Duration, wg *sync.WaitGroup) {

	for req := range requests {
		reporter := &SamplingProgressReporter{
			Progress: progress,
			Id:       req.Filename,
			Interval: interval,
		}
		Download(req.Url, req.Filename, reporter)
		wg.Done()
	}

}

func Download(url string, path string, reporter ProgressReporter) {
	reporter.ReportStart(time.Now())
	defer reporter.ReportDone()

	resp, e := http.Get(url)
	if e != nil {
		reporter.ReportError(e)
		return
	}
	defer resp.Body.Close()
	reporter.ReportTotal(resp.ContentLength)

	e = os.MkdirAll(filepath.Dir(path), os.ModePerm)
	if e != nil {
		reporter.ReportError(e)
		return
	}
	file, e := os.Create(path)
	if e != nil {
		reporter.ReportError(e)
		return
	}
	defer file.Close()

	readCounter := &ReadCounter{
		Reporter: reporter,
	}
	if _, e := io.Copy(file, io.TeeReader(resp.Body, readCounter)); e != nil {
		reporter.ReportError(e)
		return
	}
}

type ReadCounter struct {
	Current int64
	Reporter ProgressReporter
}

func (r *ReadCounter) Write(b []byte) (int, error) {
	n := len(b)
	r.Current += int64(n)
	r.Reporter.ReportProgress(r.Current)
	return n, nil
}