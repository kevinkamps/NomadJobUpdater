# Nomad Job Updater
This project basically converts a hcl file to json and posts it to the nomad managers to update/add a job. It's very useful when doing
some automatic deployment through a CI system like Jenkins / Gitlab / TeamCity / Travis CI / Bamboo / etc...

## CLI Options
```
Usage:
  -job-hcl-file string
        Path to the job hch file (default "nomad-job.hcl")
  -nomad-url string
        Parse url (default "http://127.0.0.1:4646")
  -version
        Prints the version of the application and exits
```

## Usage
### Running as Docker container: Showing help
```bash
docker run -d \
    --name=nomad-job-updater \
    kevinkamps/nomad-job-updater:latest \
      nomad_job_updater -help
```
### Running as Docker container: Updating / adding a job to your Nomad cluster
```bash
docker run -d \
    --name=nomad-job-updater \
    kevinkamps/nomad-job-updater:latest \
      nomad_job_updater -nomad-url=https://nomad.domain.com -job-hcl-file=nomad-job.hcl
```
## License

[GPL-3.0](https://choosealicense.com/licenses/gpl-3.0/)
