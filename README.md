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
