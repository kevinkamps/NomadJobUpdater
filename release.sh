#!/usr/bin/env bash

echo -e "Version: "
read VERSION

export GO111MODULE=on

mkdir bin
rm -rf bin/*

echo "$(date): Building for Linux"
GOOS=linux GOARCH=386 go build -ldflags "-X main.version=$VERSION" -o ./bin/nomad_job_updater_linux_386 ./main.go
GOOS=linux GOARCH=amd64 go build -ldflags "-X main.version=$VERSION" -o ./bin/nomad_job_updater_linux_amd64 ./main.go

echo "$(date): Building container"
docker build -t kevinkamps/nomad-job-updater:$VERSION .
docker tag kevinkamps/nomad-job-updater:$VERSION kevinkamps/nomad-job-updater:latest


echo "$(date): Pushing to docker"
echo -e "Want to push $VERSION and latest to Docker (y/n)"
read CHOICE
if [ $CHOICE = y ]
	then
	docker push kevinkamps/nomad-job-updater:$VERSION
	docker push kevinkamps/nomad-job-updater:latest
fi


echo "$(date): Done"