FROM ruby:2.6-alpine3.11

RUN gem install resque

EXPOSE 3000

CMD resque-web -p 3000 -r redis -F
