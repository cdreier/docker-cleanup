based on [this gist](https://gist.github.com/bastman/5b57ddb3c11942094f8d0a97d461b430) i reated a small docker image to clean up docker hosts

## usage

run it with with docker.sock as volume

`docker run -v /var/run/docker.sock:/var/run/docker.sock drailing/docker-cleanup:latest`

i write logs to a local file, a webserver is running on port 8080 to have easy access

`docker run -p 8080:8080 -v /var/run/docker.sock:/var/run/docker.sock drailing/docker-cleanup:latest`

or with docker-compose

```
version: '3'

services:
  app: 
    image: drailing/docker-cleanup:latest
    # restart: always
    ports:
      - 8080:8080
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock
```

## CLI

TODO