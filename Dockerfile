FROM ruby:3.1.2

WORKDIR /usr/src/app
COPY . .

RUN gem install bundler -v 2.3.7
RUN bundle config set without development
RUN bundle

CMD ["bundle", "exec", "rails", "server", "-b", "0.0.0.0"]
