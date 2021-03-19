FROM centos:centos7

RUN yum -y update; yum clean all
RUN yum -y install epel-release; yum -y install redis; yum clean all

EXPOSE 6379
CMD ["redis-server", "--protected-mode", "no"]
