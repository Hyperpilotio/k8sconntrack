# Set the base image
FROM ubuntu:16.04

# Set the file maintainer
MAINTAINER swh <swh@hsiang.io>

RUN apt-get update && \
    apt-get install iptables -y

ADD ./_output/conntracker /bin/conntracker

ENTRYPOINT ["/bin/conntracker"]
