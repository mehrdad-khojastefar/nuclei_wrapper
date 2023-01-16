package subfinder

import (
	"bytes"
	"io"
	"strings"

	"github.com/google/uuid"
	"github.com/projectdiscovery/subfinder/v2/pkg/runner"
)

type Runner struct {
	Id       string          `json:"id"` // GUID
	Name     string          `json:"name"`
	Options  *runner.Options `json:"options"`
	Instance *runner.Runner
}

// Creates a new Runner
func NewRunner(name string, opt *runner.Options) (*Runner, error) {
	// default values
	opt.JSON = false
	opt.Silent = true
	opt.Verbose = false

	r := &Runner{
		Id:      uuid.New().String(),
		Name:    name,
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
	buf := bytes.Buffer{}
	err := r.Instance.EnumerateSingleDomain("iran.ir", []io.Writer{&buf})
	if err != nil {
		return nil, err
	}

	data, err := io.ReadAll(&buf)
	if err != nil {
		return nil, err
	}
	// creating the list from output
	arr := strings.Split(string(data), "\n")
	return arr, nil
}
