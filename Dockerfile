# create a golang base image for ci/cd
FROM golang:1.23.0

# install git
RUN apt-get update && apt-get install -y git

# install buf for protobuf




