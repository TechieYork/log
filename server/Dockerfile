FROM centos:latest
MAINTAINER Techie York "york528@yeah.net"
ENV REFRESHED_AT 2017-03-29

ADD ./admin/ /data/soft/DarkMetrix/log/server/admin
ADD ./conf/ /data/soft/DarkMetrix/log/server/conf
ADD ./bin/ /data/soft/DarkMetrix/log/server/bin

VOLUME ["/tmp"]
WORKDIR "/data/soft/DarkMetrix/log/server/bin/"

ENTRYPOINT ["./dm_log_server"]
CMD ["-h"]
