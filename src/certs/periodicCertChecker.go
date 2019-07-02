package certs

import (
	"time"
	"path/filepath"

	"k8s.io/klog"
)

type PeriodicCertChecker struct {
	period time.Duration
	includeCertGlobs []string
	excludeCertGlobs []string
}

func NewCertChecker(period time.Duration, includeCertGlobs []string, excludeCertGlobs []string) *PeriodicCertChecker {
	return &PeriodicCertChecker{
		period : period,
		includeCertGlobs : includeCertGlobs,
		excludeCertGlobs : excludeCertGlobs,
	}
}

func (p *PeriodicCertChecker) StartChecking() {

	periodChannel := time.Tick(p.period)

	for {
		klog.Info("Begin periodic check")

		for _, match := range p.getMatches() {

			if !p.includeFile(match) {
				continue
			}

			klog.Infof("Publishing metrics for %v", match)
		}

		<-periodChannel
	}
}

func (p *PeriodicCertChecker) getMatches() []string {
	ret := make([]string, 0)

	for _, includeGlob := range p.includeCertGlobs {

		matches, err := filepath.Glob(includeGlob)

		if err != nil {
			klog.Errorf("Glob failed on %v: %v", includeGlob, err)
			continue
		}

		ret = append(ret, matches...)
	}

	return ret
}

func (p *PeriodicCertChecker) includeFile(file string) bool {

	for _, excludeGlob := range p.excludeCertGlobs {
		exclude, err := filepath.Match(excludeGlob, file)

		if err != nil {
			klog.Errorf("Match failed on %v,%v: %v", excludeGlob, file, err)
			return false
		}

		if exclude {
			return false
		}
	}

	return true
}