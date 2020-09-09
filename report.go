package godl

import (
	"time"
)

type Progress struct {
	Id      string
	Current int64
	Total   int64
	Speed 	float64
	Elapsed time.Duration
	Error   error
}

type ProgressReporter interface {
	ReportProgress(progress int64)
	ReportTotal(total int64)
	ReportStart(start time.Time)
	ReportError(err error)
	ReportDone()
}

type SamplingProgressReporter struct {
	Current  int64
	Total    int64
	Progress chan<- Progress
	Id       string
	Interval time.Duration
	Start    time.Time
	Last     time.Time
	Error    error
}

func (r *SamplingProgressReporter) ReportProgress(progress int64) {
	r.Current = progress
	if time.Since(r.Last) > r.Interval {
		r.Progress <- Progress{
			Id:      r.Id,
			Current: r.Current,
			Total:   r.Total,
			Elapsed: time.Since(r.Start),
		}
		r.Last = time.Now()
	}
}

func (r *SamplingProgressReporter) ReportTotal(total int64) {
	r.Total = total
}

func (r *SamplingProgressReporter) ReportStart(start time.Time) {
	r.Start = start
}

func (r *SamplingProgressReporter) ReportError(err error) {
	r.Error = err
	r.Progress <- Progress{
		Id:    r.Id,
		Error: err,
	}
}

func (r *SamplingProgressReporter) ReportDone() {
	r.Progress <- Progress{
		Id:      r.Id,
		Current: r.Current,
		Total:   r.Total,
		Elapsed: time.Since(r.Start),
		Error:   r.Error,
	}
}
