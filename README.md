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

The reason why we can't just get accounts and update on the fly is that multiple concurrent access to the database can still access the same stale data. So we need to make sure incoming database requests wait for a cell to finish updating. We're doing that by adding FORCE UPDATE to the sql query. 