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

## HCL Templates
Details about this specification can be found at: https://github.com/hashicorp/hcl

### Variables
We have added support for variables in the template. We only support replacement of environment variables. You can use `$` followed by a 
variable name `$variable_name` that would be replaced the environment variable value of `variable_name`. We do not support the convention 
of `${var}` which is widely used in the industry because that convention is used by nomad (see: https://www.nomadproject.io/docs/runtime/interpolation.html#interpreted_node_vars)

## License

[GPL-3.0](https://choosealicense.com/licenses/gpl-3.0/)
