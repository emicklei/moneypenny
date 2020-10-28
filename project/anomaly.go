package project

import "fmt"

type AnomalyDetector interface {
	IsAnomaly(stats *ProjectStats) bool
	String() string
}

// https://github.com/SundaySky/cost-anomaly-detector
type SundaySky struct {
	relativeThreshold float64
	stddevThreshold   float64
	absoluteThreshold float64
	windowDays        int
}

var BestSundaySky = SundaySky{1.25, 2.0, 10, 30} // article uses 3.5 for stddev

func (s SundaySky) String() string {
	return fmt.Sprintf("sundaysky{relativeThreshold=%.2f,stddevThreshold=%.2f,absoluteThreshold=%.2f,windowDays=%d}", s.relativeThreshold, s.stddevThreshold, s.absoluteThreshold, s.windowDays)
}

func (s SundaySky) IsAnomaly(stats *ProjectStats) bool {
	if len(stats.Daily) < 1 {
		return false
	}
	// on day end
	charge := stats.Daily[0].Charges
	if charge <= s.relativeThreshold*stats.Mean {
		return false
	}
	if charge <= s.stddevThreshold*stats.StandardDeviation {
		return false
	}
	if charge <= s.absoluteThreshold {
		return false
	}
	return true
}
