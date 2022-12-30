#To convert this file from a normal build file to a multistage build file, we used the AS keyword and the stage name to the end of the FROM keyword
#id `AS [name] to specify the current stage`

#The Build Stage
#This specifies the version of the golang image to use
FROM golang:1.19.4-alpine3.17 AS build_stage

#The working directory inside the image, we kept it simple, used `/app`
WORKDIR /app

#This copy command copies all the golang files into the working directory
COPY . .

#This RUN command runs the command in front of it, building an executable file
RUN go build -o main main.go

# ----- We no longer need to download the migrate binary because we are now running DB migrations direcly in the golang code
# This command should add curl
# RUN apk add curl

# #This RUN command downloads the migrate tool so we can call migrate up to create the tables for the database
# RUN  curl -L https://github.com/golang-migrate/migrate/releases/download/v4.12.2/migrate.linux-amd64.tar.gz | tar xvz


#The Run State
FROM alpine:3.17 AS run_stage
WORKDIR /app

#Copying the executable file from the build stage
#Note: The order of execution of these commands matter. For one, you copy start.sh before copying db/migration
COPY --from=build_stage /app/main .

# No need to copy the migrate binary file too as we are now running DB migrations direcly in the golang code
# COPY --from=build_stage /app/migrate.linux-amd64 ./migrate

COPY app.env .
COPY start.sh .
COPY wait-for.sh .
COPY db/migration ./db/migration


#This specifies the port that the application will be listening on
EXPOSE 8080

#This specifies the default command to run when the image starts, CMD is an array of command line arguments
CMD [ "/app/main" ]
ENTRYPOINT [ "/app/start.sh" ]