FROM ruby:2.6-alpine3.11

RUN gem install resque


RUN mkdir /resque-data

COPY job.rb /resque-data
COPY add_job.rb /resque-data

CMD ["/bin/sh"]
