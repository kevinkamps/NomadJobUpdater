FROM busybox:1.31-glibc
LABEL maintainer="Kevin Kamps"
LABEL github="https://github.com/kevinkamps/NomadJobUpdater"
LABEL license="GPL-3.0"

COPY bin/nomad_job_updater_linux_amd64 /bin/nomad_job_updater