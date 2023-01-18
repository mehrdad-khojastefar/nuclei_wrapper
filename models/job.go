package models

type JobStatus int

const (
	JOBSTATUS_FINISHED JobStatus = iota
	JOBSTATUS_ONGOING
	JOBSTATUS_ERROR
)

type Result struct {
	Domain       string `json:"domain" bson:"domain"`
	Ip           string `json:"ip" bson:"ip"`
	NucleiResult string `json:"result,omitempty" bson:"result,omitempty"`
}

type Job struct {
	Id         string    `json:"scan_id" bson:"_id"`
	Status     JobStatus `json:"status" bson:"status"`
	Date       string    `json:"date" bson:"date"`
	Domain     string    `json:"domain" bson:"domain"`
	SubDomains []Result  `json:"subdomains,omitempty" bson:"subdomains,omitempty"`
}
