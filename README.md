# Bank service

## Features

- [x] Process transactions between users
- [x] Store history of transactions in database
- [x] Saving state of RabbitMQ
- [x] Saving state of Postgres
- [ ] Publishing messages to RabbitMQ

## How to use

1. Run *make compose-app* at the root of the project to run application with all dependencies in docker-compose

## To send messages to RabbitMQ use the following instructions:

1. Run *make compose-app* at the root of the project to run application
2. In your browser go to http://localhost:15672
3. Login with username *rmuser* and password *rmpassword*
4. Go to tab *Queues and streams*
5. Select *transaction-queue*
6. Expand *Publish message* at the bottom of the page 
7. Fill in the form with example.json file 
8. Click on *Publish message* button
