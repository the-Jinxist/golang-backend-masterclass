#This specifies the version of the golang image to use
FROM golang:1.18.8-alpine3.16

#The working directory inside the image, we kept it simple, used `/app`
WORKDIR /app

#This copy command copies all the golang files into the working directory
COPY . .

#This RUN command runs the command in front of it, building an executable file
RUN go build -o main main.go

#This specifies the port that the application will be listening on
EXPOSE 8080

#This specifies the default command to run when the image starts, CMD is an array of command line arguments
CMD [ "/app/main" ]