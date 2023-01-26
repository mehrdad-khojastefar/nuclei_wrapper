package main

import (
	"bytes"
	"io"
	"strings"

	"github.com/google/uuid"
	"github.com/projectdiscovery/subfinder/v2/pkg/runner"
)

type Runner struct {
	Id       string          `json:"id"` // GUID
	JobId    string          `json:"job_id"`
	Options  *runner.Options `json:"options"`
	Instance *runner.Runner
}

// Creates a new Runner
func NewRunner(jId string, opt *runner.Options) (*Runner, error) {
	// default values
	opt.JSON = false
	opt.Silent = true
	opt.Verbose = false
	opt.HostIP = true
	opt.RemoveWildcard = true

	r := &Runner{
		Id:      uuid.New().String(),
		JobId:   jId,
		Options: opt,
	}
	instance, err := runner.NewRunner(r.Options)
	if err != nil {
		return &Runner{}, err
	}
	r.Instance = instance
	return r, nil
}

// runs the domain enumeration and returns the results
func (r *Runner) GetSubdomainArray(subdomain string) ([]string, error) {
	buf := &bytes.Buffer{}
	err := r.Instance.EnumerateSingleDomain(subdomain, []io.Writer{buf})
	if err != nil {
		return nil, err
	}
	data, err := io.ReadAll(buf)
	if err != nil {
		return nil, err
	}
	// creating the list from output
	// fmt.Println(string(data))
	arr := strings.Split(string(data), "\n")
	return arr, nil
}
