## Service for xlsx files generation.

### Commands:

#### Docker build

````shell
docker build -t xsls-generator:latest .
````

#### Docker run

````shell
docker run -d -t -i -p 8080:8080 --name xsls-generator xsls-generator:latest
````