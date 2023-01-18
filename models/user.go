package models

import (
	"time"

	"github.com/google/uuid"
	"hamravesh.ir/mehrdad-khojastefar/database"
)

type User struct {
	Id       string `json:"id" bson:"_id"`
	Username string `json:"username" bson:"username"`
	Password string `json:"password" bson:"password"`
	Jobs     []Job  `json:"jobs,omitempty" bson:"jobs,omitempty"`
}

func (u *User) GetJobs() []Job {
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
	u.Jobs = append(u.Jobs, j)
	// update the database with new user
	err := database.Db.UpdateUser(u)
	if err != nil {
		return "", err
	}
	return j.Id, nil
}
