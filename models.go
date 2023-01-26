package main

import (
	"encoding/json"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/projectdiscovery/nuclei/v2/pkg/output"
	"github.com/projectdiscovery/subfinder/v2/pkg/resolve"
	"github.com/projectdiscovery/subfinder/v2/pkg/runner"
)

type JobStatus int

const (
	JOBSTATUS_ONGOING JobStatus = iota
	JOBSTATUS_FINISHED
	JOBSTATUS_ERROR
)

func (s JobStatus) String() string {
	switch s {
	case JOBSTATUS_ERROR:
		return "ERROR"
	case JOBSTATUS_ONGOING:
		return "ONGOING"
	case JOBSTATUS_FINISHED:
		return "FINISHED"
	}
	return "INVALID_STATUS"
}

type Result struct {
	Domain   string `json:"domain" bson:"domain"`
	Ip       string `json:"ip" bson:"ip"`
	Resolver string `json:"resolver" bson:"resolver"`
}

type Job struct {
	Id           string                `json:"scan_id" bson:"_id"`
	Status       JobStatus             `json:"status" bson:"status"`
	Date         string                `json:"date" bson:"date"`
	Domain       string                `json:"domain" bson:"domain"`
	SubDomains   []*Result             `json:"subdomains,omitempty" bson:"subdomains,omitempty"`
	NucleiResult []*output.ResultEvent `json:"nuclei_result,omitempty" bson:"nuclei_result,omitempty"`
}

var wg sync.WaitGroup

func (j *Job) StartNewJob(user *User) error {
	go func() error {
		// simple quee logic
		wg.Wait()
		wg.Add(1)
		defer wg.Done()

		j.Status = JOBSTATUS_ONGOING
		err := Db.UpdateJob(j.Id, user)
		if err != nil {
			return err
		}
		// first get the subdomains
		instance, err := NewRunner(j.Id, &runner.Options{
			Threads:            10,
			Timeout:            30,
			MaxEnumerationTime: 10,
			Resolvers:          resolve.DefaultResolvers,
		})
		if err != nil {
			return err
		}
		subdomains, err := instance.GetSubdomainArray(j.Domain)
		if err != nil {
			return err
		}

		domains := make([]string, 0, len(subdomains))
		for _, v := range subdomains {
			d := strings.Split(v, ",")
			if v == "" || v == " " {
				continue
			}

			j.SubDomains = append(j.SubDomains, &Result{
				Domain:   d[0],
				Ip:       d[1],
				Resolver: d[2],
			})
			domains = append(domains, d[0])
		}
		if err != nil {
			return err
		}
		err = Db.UpdateJob(j.Id, user)
		if err != nil {
			return err
		}
		nucleiResult, err := StartNuclei(domains)
		if err != nil {
			if len(nucleiResult) != 0 {
				j.NucleiResult = nucleiResult
				j.Status = JOBSTATUS_FINISHED
				_ = Db.UpdateJob(j.Id, user)
			}
			return err
		}
		j.NucleiResult = nucleiResult
		j.Status = JOBSTATUS_FINISHED

		err = Db.UpdateJob(j.Id, user)
		if err != nil {
			return err
		}
		return nil
	}()
	return nil
}

type User struct {
	Id       string `json:"id" bson:"_id"`
	Username string `json:"username" bson:"username"`
	Password string `json:"password" bson:"password"`
	Jobs     []*Job `json:"jobs,omitempty" bson:"jobs,omitempty"`
}

func (u User) MarshalJSON() ([]byte, error) {
	type user User // prevent recursion
	x := user(u)
	for _, j := range x.Jobs {
		j.SubDomains = nil
	}
	return json.Marshal(x)
}

func (u *User) GetJobs() []*Job {
	// TODO: pagination
	return u.Jobs
}

func (u *User) AddJob(domain string) (string, error) {
	j := Job{
		// set the job id
		Id: uuid.New().String(),
		// set the job date
		Date: time.Now().UTC().String(),
		// set the status for the job
		Status: JOBSTATUS_FINISHED,
		Domain: domain,
	}
	// TODO:add the job to the quee of all the jobs
	u.Jobs = append(u.Jobs, &j)
	// update the database with new user
	err := Db.UpdateUser(u)
	if err != nil {
		return "", err
	}
	return j.Id, nil
}

func (u *User) HasJob(jobId string) (bool, *Job) {
	for _, v := range u.Jobs {
		if v.Id == jobId {
			return true, v
		}
	}
	return false, nil
}

func (u *User) ConvertToMap() (map[string]interface{}, error) {
	userByte, err := json.Marshal(u)
	if err != nil {
		return nil, err
	}
	var userJson map[string]interface{}
	err = json.Unmarshal(userByte, &userJson)
	if err != nil {
		return nil, err
	}
	return userJson, nil
}
