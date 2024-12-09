# create a golang base image for ci/cd
FROM golang:1.23.0

# install git
RUN apt-get update && apt-get install -y git

# install make
RUN apt-get install -y make

# install curl
RUN apt-get install -y curl

# install jq
RUN apt-get install -y jq

# install buf for protobuf


# install docker
RUN apt-get install -y docker.io


