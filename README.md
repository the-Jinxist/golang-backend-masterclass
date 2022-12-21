# Another Learning Project

This project is supposed to help me learn how to use Golang + Postgres + Kubernetes + gRPC by building a banking
application lol. 

DAY 1:

Used dbdiagrams.io to design a database schema. Will attach PDF and SQL commands in codebase

DAY 2:

Learnt how to use Docker, download images, create containers, execute some sql commands via docker in terminal. Learnt
how to use TablePlus to run commands is postgres environment

DAY 3:

Learnt database migration using the golang/migrate tool. Learnt about up-down migration which is apparently best practice for database migration. The up script is run when we want to make a forward change to the schema and the down script is run when we want to reverse a change to the schema

DAY 4:

Will put all progress till now here. Used the sqlc pacjage to generate sql golang code using sql files in the db/queries folder. These generated golang sql operation code can be found in the db/sqlc folder, the end with .sql.co
Also worked on writing unit tests for these sql operations. They can be found in the db/sqlc files.

Day 5:

Learning about Database Transaction: A unit of work with multiple database operations. In our case
    - Creating a transfer record for a transfer
    - Create entry for account 1 with -money sent
    - Create entry for account 2 with +money recieved
    - Subtract money from account 1
    - Add money to account 2

    Db Transactions:

    - Are needed for reliability
    - Are needed for isolation between multiple programs accessing data concurrently

    Db Transactions must follow the ACID rule. 

    - Atomicity: All database operations must be successful or fail and the db is unchanged
    - Consistency: DB state must be valid after the transaction is complete.
    - Isolation: Concurrent transactions must not affect each other
    - Durability: All data written be a successful transaction must be recorded in persistent storage.

Day 6: 

Learnt a clean way to implement database transactions. Normally, you start like a transaction session, get the query object from it using the New() method we created, use the query object to make multiple operations while keeping track fo the errors that could happen in each operation. if any error is found pass it back and rollback the transaction, else, commit the transaction.

Day 7: 

The reason why we can't just get accounts and update on the fly is that multiple concurrent access to the database can still access the same stale data. So we need to make sure incoming database requests wait for a cell to finish updating. We're doing that by adding FOR UPDATE to the sql query. 

Day 8: 

Simulated a database deadlock. Worked on debugging that using a transaction key and name for each transaction using context.WithValue() and logging those values. We fiund two ways to avoid deadlocks in this case: removing the foreign key constraints(a bad solution because it reduces the validity of the database) and adding the FOR NO KEY UPDATE to the sql query. We also updated the update account sql command to also change the add an amount to the balance of an account

Day 9:

We learnt that deadlocks can still happen. The best way to avoid them is to make sure that you order your query very well. We also studied isolation levels within query transactions in postgres. Highlighted the I in ACID. We have to make sure that transactions running concurrently don't affect each other. This can result in multiple forms.
- Dirty Read: A transaction reads data written by another concurrent *uncommited* transaction. This is very bad because we don't even know it the uncommited transactionn will be actually commited or rolled back.
- Non-repeatable read: A transaction reads a row twice and sees different values because it has been modified by other committed transaction.
- Phantom read: A transaction re-executes a query to find rows that meet a certain condition and sees a different set of rows due to changes by other commited transaction
- Serialization Anomaly: The result of a group of concurrent commited transactions is impossible to achieve if we try to run them sequentially in any order without overlapping

To beat this 4 isolation leves were created:
- Read Uncomitted: Transactions can read uncommited changes to the database
- Read Comitted: Transactions can only read commited changes to the database
- Repeatable Read: Same read query always return the same result
- Serializable: Can achieve same result if transactions are executed in sequential order

Day 10:
    In order to see how the isolation levels work with read phenomena, I had to get the docker image for mysql and create a container, connect to table plus and run commands in the docker container using docker exec. It was something. 

Day 11;
Finished learning about isolation levels in MySql and PostGres.
- MySql:
    - Read Uncommitted: Doesn't prevent any of the transaction phenomena.
    - Read committed: Only prevents dirty read.
    - Repeatable Read: Prevents dirty read, Non-repeatable read and phantom read.
    - Serializable: Prevents every transaction phenomena.

    - Uses locking mechanisms for stoping concurrent reads when a share lock is gotten.
    - The default transaction isolation level: is repeatable read.

- Postgres: 
    - Read Uncommitted: Behaves the same way as read commited. Only prevents dirty read.
    - Read committed: Only prevents dirty read.
    - Repeatable Read: Prevents dirty read, Non-repeatable read and phantom read.
    - Serializable: Prevents every transaction phenomena.

    *So Postgres basically has 3 isolation levels*

    - Uses dependency detection to detect when a share lock is gotten
    - The default transaction isolation level: is read uncommitted

NB: Using a high transaction isolation level might lead to deadlocks so you have to implement retry mechanisms. Also make sure to read documentation to see how database engines implement
    these transaction isolation levels.

Setting up Github Actions to run automated tests:
1. Workflow: Automated procedure. Made up of 1+ jobs. Triggered by events, scheduled or manually. To create a workflow, added a .yml file to your repository.
2. Runner: A server that listens for available jobs, runs one job at a time, we can use a github-hosted runner or a runner of our choosing. Reports progress and logs and results
   to the Github UI 
3. Jobs: A set of steps that execute on the same runner. Normal jobs run in parallel. Dependent jobs run serially.
4. Steps: an individual task that run serially within a job. contains 1+ actions
5. Actions: A standalone command. Run serially within a step. Can be reused.

Day 12:

Finsihed up Github Actions. Main blockers what using the service container for postgres, adding the necessary evironment variables. More details can be found (here)[https://docs.github.com/en/actions/using-containerized-services/creating-postgresql-service-containers], and running migrate tools with the right CLI. Links can be found (here)[https://github.com/golang-migrate/migrate#cli-usage] also.

Day 13:

Started work on using the (Gin)[https://github.com/gin-gonic/gin] library to create REST Api endpoints.

Day 14:

Created routes for creating and getting accounts.

Day 15:

Created the route for getting list of accounts with pagination variables in the URL Queries. Also edited sqlc.yaml file to `emit_empty_queries: true` to make sure it returns an empty list when there are no more accounts to query. Finished up HTTP API implementation wth golang. Finished up initial REST API learnings with golang

Made 3 endpoints:
- GET with `/account/:id` to get one unique account
- POST with `/accounts` to create one account.
 - With a payload like:
   ```
    {
        "owner": "Favour",
        "currency": ""
    }
    
   ```
- GET with `/accounts` to GET  a list of accounts. Requires `page_size` and `page_id` parameters to work.

*Local host URL, of course `localhost:8080`*

Day 16:

Learnt how to load environment variables with Viper package. Moved away from hardcoding environment variables to adding them to an app.env file

Day 17:

Learning how to mock DBs for testing. Why do we mock DBs

Why do we mock database(not to make fun of them, I promise🌚):

1. Independent Tests: Isolate tests to avoid conflicts
2. Faster tests: Since they don't use the actual DB
3. 100% Coverage: Mock DBs can be used to test unexpected erros and results, which cannot be done using an actual DB

How to mock:

1. Use fake DM - Memory: Implement a fake version of DB: store data in memory. However this requires us to write a lot more code
2. Use DB Stubs - Gomock: Generate and build stubs that return hard-coded values

`var _ Querier = (*Queries)(nil)` in the querier.go file makes sure that the `Queries` struct implements all the methods that are spelt out in
the `Querier` interface. We then made the `Store` struct into an `interface`, and we used the `gomock` package to mock the Store struct.

*Using mockgen*: Mockgen has two ways of mocking interfaces/structs. `reflect` and `source` mode. `source` mode will get very much complicated if we have other code
imported from other files, `reflect` just makes use of the package of the file and the name of the interface. So we choose to use `reflect` instead.

running the mock command goes like: `mockgen {[module copied from top of go.mod file]/[path]/[to]/[interface]/[you]/[want]/[to]/[mock]} {name of the interface}`

In order to make it work, I add to run `go mod tidy -compat=1.17` and `go mod vendor`. For any missing values that don't have their respective folder in the generated vendor foler, you can just do a blank import in any file and run `go mod tidy` and `go mod vendor` again

Day 18:

Learnt a lesson. Make sure the struct you're returning in `recorder.Body([struct])` is the same as the struct returned in `ctx.Json(code, [struct]`). Unless your test will not run o lmao. Main issues:
- Creating a list of anonymous structs kinda caught me off guard, anonymous structs containing anonymous functions too lol
- made sure the forloop uses all the vales from each item in the list of test cases.
- Created a new main_test.go file so we can remove the verbose logs that gin keeps giving us

Day 19:

We created the endpoint for creating transfers and wrote up a custom validator for Gin to check if the currency used in every transaction is supported by our simplebank. We also made sure to check
for discrepancies in currencies for different accounts. Maybe I will whip up a transfer rate someday sha. Something to think about

Day 20:

We created a new migration with the new postgresql code generated from the dbdiagram.io model. Using the command `migrate create -ext sql -dir {directory-to-your-migration-files} -seq {name-of-migration}`. We're adding a new table called users and working on the migration

Day 21:

Finished writing up the SQL commands for the migrate up and migrate down for the new `users` table.

Day 22: Learnt how to handle DB errors in golang. Also updated our code because we can't use any random owner for creating account anymore as per our created constraints.

Day 23: Hashing password with bcrypt. Created util functions to hash and compare passwords, wrote tests for it too.

Day 24: Added API for create and get user. Tweaked response to not return hashed password.

Day 25: Used a custom gomock matcher to make sure we pass the arguments with the right values to the endpoint.

Day 26: We learnt about how Paseto is better than JWT for authentication.

Day 27: Created the JWT maker file. This file implements our abstract token creator interface.

Day 28: Learned how to verify tokens with JWT and Paseto. Wrote test to account for the way each implementation works

Day 29: Created a new endpoint for login. Applied and interchanged between Paseto and JWT tokem validations

Day 30: Added Auth Middle ware to routes registration and also added implementation for adding and checking token to the header for each request. Wrote tests for checking the implementation for the 
auth middleware and somehow understood more the flow of testing http requests using recorder and the http package

Day 31: Adding authorization rules to database endpoints

Day 32: Added authorization test changes to the transfers and account test files. Started working on pushing a release version into the wild.
        Using a multistage docker image file to make this happen.
        We called `docker build -t simplebank:latest .`

Day 33: Learning how to connect two stand-alone container. 
    1. First method was finding a way to manually connect these containers using the IP address of the other container. This IP address could be foubd
    using the `docker container inspect <container-id or image name>`
    2. Second method, the preferrable one was the Network method. Creating a network in which both containers would be connected
    Networks can be created using `docker network create <network-name>`
    The we add containers to a network using `docker network connect <network-name> <container-name or container id>`
    The we ran the image for the new simple bank project we created using the multistage on the same network using the following complicated command: 
        `docker run --name simplebank --network bank-network -p 8080:8080 -e GIN_MODE=release -e DB_SOURCE="postgresql://root:secret@postgres-learning1:5432/simple_bank?sslmode=disable" simplebank:latest`

Day 34: Learning how to write a docker-compose file and define a startup flow for services. To run the docker compose file, use. `docker compose up`. Created a `startup.sh` file. We ran `chmod +x start.sh` to change the mode of the file to an executable. We used #!bin/sh to run the startup.sh file because we're using alpine12, as the bash shell is not available.

Dy 35: We are working on running migrations after the postgres container is created. To do this, we learn't about the wait-for sh file. Not to make .sh files executable run `chmod +x [sh-file-name].sh`

Day 36: We've created a free-tier AWS account and I connected my payday account to it. Worked like a charm. We decided that the flow to deploy our containers to production via AWS will be done when we push to main, so we used it in a github action. We found github actions for using AWS ECR, such as logging in and adding the credentials. at https://github.com/marketplace/actions/amazon-ecr-login-action-for-github-actions. We decided to create an IAM user to safegaurd our root user credentials. We'll use this IAM user's credential to make the necessary pushes and pulls needed.

We used the Secrets section on the left-hand side in the Github Repo tab to add environment variables to our actions. Apparently, there are two types of secrets, Environment Secrets for branches like staging, main, e.tc. and Repository secrets for the entire repository.

Day 37: The deployment didn't work until we added the permission `AmazonEC2ContainerRegistryFullAccess` to the iam user. Two solid days of work gone. We worked with using Amazon RDS to create a production database. We chose the `Standard Create` option. The DB engine we used in Postgres, the version `12.12-R1` and the `Free-tier` template. We set our Master username to `root` and we checked the field to make RDS autogenerate a password for us. We disabled `Enable storage auto scaling`. We allowed public access to our database and setup a security policy for it. Named it `access-postgres-anywhere` and gave no preference for `Availability Zone`

In Big `Additional Configuration`, We added a name to the initial database that Amazon RDS will create for us: `simple_bank`. We opened the details for the newly-created simple bank RDS database. Clicked the link for the VPC security group. Tapped the id for our security group, And we edited the inbound rules to not only allow access from my ip address. We should not do this on a production DB.

We connected TablePlus to our online database, using the Master password and username we've generated plus the endpoint for the database given to us by AWS RDS. After this initial connection the database is very empty, so we have to run the migrate up command from the makefile but we'll update the username and password parts i.e `root:secret`. The host was the link/endpoint to the database generated by AWS.

All in all, we learnt how to create and set up a database instance on AWS RDS.

Moving on, we worked with Amazon's Secrets Manager to store secrets for the database. We chose `Other type of secret` secret type. We added the key value pairs for the variables in `app.env`. We replaced the DB_SOURCE with the new one we used to run migrate up with. Using root:[generated-password]@[generated-databse-endpoint]:[port]/[database-name]

Day 36: We've started using Secrets Manager. We used the openssl command `openssl rand -hex 64 | head -c 32` to generate a 32-length random string for our token symmetric key. After creating our secret, We moved on to updating our github workflow, to do this we installed the aws cli. After installing the cli, we ran `aws configure` to set up our credentials for accessing AWS. We opened our `github-ci` user in IAM service in AWS, went to the security credentials tab, then found our access key id and secret key values there.

We then use the command `aws secretsmanager get-secret-value --secret-id simplebank`. We had to give the user the permission to access to secret values. We did this by going the User groups page and adding a new permission policy, we found erm  `SecretsManagerReadWrite` policy and attached it to the user group. To get just specifically the secret string, we use `aws secretsmanager get-secret-value --secret-id simplebank --query SecretString` to get the output in json format, we call `aws secretsmanager get-secret-value --secret-id simplebank --query SecretString --output text`

To convert the output generated from the above command into an .env file-ish format, we turned to the jq library here: https://stedolan.github.io/jq/. We installed jq using `arch -arm64 brew install jq`. We chained the command to get-secret-values from aws with the pq command `pq to-entries` to get an array of key-value pairs for each environemnt variable. Command lokked like this `aws secretsmanager get-secret-value --secret-id simplebank --query SecretString --output text | jq 'to_entries`. We then used the string interpolation feature of the jq library to format the values in the env style. The command looked like this. `aws secretsmanager get-secret-value --secret-id simplebank --query SecretString --output text | jq 'to_entries|map("\(.key)=\(.value)")'`.

To remove the array symbols `[]` from the output of the above, we use `aws secretsmanager get-secret-value --secret-id simplebank --query SecretString --output text | jq 'to_entries|map("\(.key)=\(.value)")|.[]'`

To remove the `""` i.e the quotes characters from each string, we need to use the -r tag when running the command, so it'll look like this omo `aws secretsmanager get-secret-value --secret-id simplebank --query SecretString --output text | jq -r 'to_entries|map("\(.key)=\(.value)")|.[]'`. The output now look exactly as we want it, next thing is to redirect this values onto our app.env file. In essence, overwriting the contents of that file. The command will look like this: `aws secretsmanager get-secret-value --secret-id simplebank --query SecretString --output text | jq -r 'to_entries|map("\(.key)=\(.value)")|.[]' > app.env`.

We now added this super command to our github deploy workflow. So the code that will be deployed will use the production environment variables stored in our AWS secret. 

We tried to pull the newly-created image on our ECR but we don't have the necessary credentials so we turned to the command to get a login password for AWS's ECR `aws ecr get-login-password`. In order to make it work with the docker, we piped in the password, using the command as follows: `aws ecr get-login-password | docker login --username AWS --password-stdin [url of private repository without the final part i.e remove the repository name:id]`. After the login succeeds, we then pull the production image to local. using `docker pull [full private repo link]`.

So we tried running our newly-pulled local image and we had an error because our environment variables was empty. In order to rectify this, we ran `source app.env`. This command loaded all the values from our app.env file. We added this command into the start.sh file so it gets executes during startup.

Day 37: We're learning about AWS EKS and Kubernetes. K8 is a container orchestration engine for automating scaling, deployment and management of containerized applications. K8 Components are made up of a 
- Worker Node: Which always has a kubelet agent that makes sure all the containers in each k8 pods are working fine. It also has a kube-proxy which maintains network communication in the Worker node and in between k8 pods.
- The Control Plane: This runs on the master node. It's responsibility is to manage the worker nodes and pods of the cluster. The control plane consists of:
    - The API server, which is the front end of the control plane. It exposes K8 API to interact with all other components of the cluster. The persistence store etcd, which acts as K8's backing store
    for all cluster data. The scheduler which watches for newly created Pods with no assigned nodes and selectes the nodes for them to run on. The controller manager, which is a combination of several controllers such as node controller, job controller, endpoint controller, service account and token controllers, cloud controller manager components.
We mobved on the AWS Console. Went to EKS part and started work on creating a cluster. We created a service role that has permission to use the EKS cluster.

Day 38: 

We're woking on using kubectl to run k8 commands in AWS clusters? We added the kubectl cli tool using the command: `brew install kubectl`. To check if it is installed, we use `kubectl version --client`. We give kubectl access to our cluster using the following command: ` aws eks update-kubeconfig --name [cluster name] --region [cluster region]`. Our user doesn't currently have the permissions to DescribeCluster so we go to our user group, to the Permissions tab, Add a new permission, Create a new inline policy, Selected the EKS service and check all the services for EKS. To check if we can connect to our cluster via kubectl we use: `kubectl cluster-info`. On first run we had an unauthorized error. The process to solve this error is put in this link: https://aws.amazon.com/premiumsupport/knowledge-center/amazon-eks-cluster-access/. We created credentials for our root user then we updated the config file for AWS to use them.

To tell AWS cli to use the github credentials, we use the coommand `export AWS_PROFILE=github`, for default use "default" instead of "github". We created a new profile when we edited the credentials for the AWS config. We added an aws-auth.yaml file to update the users that are allowed to access the EKS. We are now learning about using K9s to work with kubernetes. We took a tour of the k9s library and it's UI to help us work with EKS clusters.

We also worked with deploying using deployment.yaml and service.yaml.

Day 39: We're working with the AWS service Route53 to buy a new domain name to host our backend application with. I couldn't buy a domain from Route53 So we couldn't get a cool domain name for our backend service and setup TLS for HTTPS support. I skipped that and moved to the automation of deploying all the k8 services when we push.
