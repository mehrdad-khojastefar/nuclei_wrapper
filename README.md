# hamravesh.ir / mehrdad-khojastefar - nuclei wrapper
the postman collection is also included in the repo 

Nuclei functionality is not stable!
having internet issues at the time of development made it hard and time consuming to test the nuclei functionality.

## Run
`docker-compose up -d`

1- first you need to create a user

2- then create a job

3- and finally start the job

4- subfinder will do its job after that nuclei starts

for the quee part Nuclei integrated it by itself.


**I found some bug in the agent part for subfinder subscraping packge and at first I vendored the modules but since it was just 1 line I manually edit it from the terminal at build time in Dockerfile**
## It was really fun to work on this task. Thanks 