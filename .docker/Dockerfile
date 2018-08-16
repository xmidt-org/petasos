FROM centos:6
MAINTAINER Jack Murdock <jack_murdock@comcast.com>

ADD ./petasos petasos
RUN mkdir -p /etc/petasos
VOLUME /etc/petasos

ENTRYPOINT ./petasos
