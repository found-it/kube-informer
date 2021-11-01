package report

import (
	"github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
)

type ReportController struct {
	report Report
	logger *logrus.Entry
}
type Report struct {
	Items []Item
}

type Item struct {
	Namespace string
	Images    []Image
}

type Image struct {
	digest     string
	tag        string
	repository string
}

func NewReportController() *ReportController {
	return &ReportController{
		logger: logrus.WithField("pkg", "report"),
		report: Report{},
	}
}

func (r *ReportController) Add(pod *v1.Pod) {
	r.logger.WithField("action", "add").Info(pod.GetName())
}

func (r *ReportController) Update(pod *v1.Pod) {
	r.logger.WithField("action", "update").Info(pod.GetName())
}

func (r *ReportController) Delete(pod *v1.Pod) {
	r.logger.WithField("action", "delete").Info(pod.GetName())
}
