FROM centos:latest
MAINTAINER Techie York "york528@yeah.net"
ENV REFRESHED_AT 2017-03-29

ADD ./admin/ /data/soft/DarkMetrix/log/agent/admin
ADD ./conf/ /data/soft/DarkMetrix/log/agent/conf
ADD ./bin/ /data/soft/DarkMetrix/log/agent/bin

VOLUME ["/var/tmp"]
WORKDIR "/data/soft/DarkMetrix/log/agent/bin/"

ENTRYPOINT ["./dm_log_agent"]
CMD ["-h"]
