FROM centos:latest
MAINTAINER Techie York "york528@yeah.net"
ENV REFRESHED_AT 2017-03-29

ADD ./test_log /data/soft/DarkMetrix/log/test/bin/

VOLUME ["/var/tmp"]
WORKDIR "/data/soft/DarkMetrix/log/test/bin/"

ENTRYPOINT ["./test_log"]
