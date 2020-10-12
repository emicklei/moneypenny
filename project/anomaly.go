package project

// https://github.com/SundaySky/cost-anomaly-detector
type SundaySky struct {
	relativeThreshold float64
	stddevThreshold   float64
	absoluteThreshold float64
}

var BestSundaySky = SundaySky{1.25, 2.0, 10} // article uses 3.5 for stddev

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
