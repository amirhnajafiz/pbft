# News feeder

Using golang framework (Gin) to create a web-application.<br />
This simple web-app uses golang gin framework to create a rest-api
application.

## Setup on local
After cloning the project run the following command:
```shell
make dev
```

And this will create the server on **localhost:8080**

## Setup on docker
You can run the application by the following steps.<br />
First build the docker image:
```shell
docker build -t newsfeed .
```

Then run it on a container:
```shell
docker run -rm -it -p 8080:8080/tcp newsfeed 
```

And then you should be able to see the app available on **localhost:8080**

## Features
- gin framework
- docker
- rest-api