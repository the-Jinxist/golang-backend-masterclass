#To convert this file from a normal build file to a multistage build file, we used the AS keyword and the stage name to the end of the FROM keyword
#id `AS [name] to specify the current stage`

#The Build Stage
#This specifies the version of the golang image to use
FROM golang:1.18.8-alpine3.16 AS build_stage

#The working directory inside the image, we kept it simple, used `/app`
WORKDIR /app

#This copy command copies all the golang files into the working directory
COPY . .

#This RUN command runs the command in front of it, building an executable file
RUN go build -o main main.go


#The Run State
FROM alpine:3.16 AS run_stage
WORKDIR /app

#Copying the executable file from the build stage
COPY --from=build_stage /app/main .

#This specifies the port that the application will be listening on
EXPOSE 8080

#This specifies the default command to run when the image starts, CMD is an array of command line arguments
CMD [ "/app/main" ]