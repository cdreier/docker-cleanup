FROM docker:18-dind

RUN mkdir /app
ADD docker-cleanup /app
WORKDIR /app

EXPOSE 8080

CMD [ "./docker-cleanup" ]