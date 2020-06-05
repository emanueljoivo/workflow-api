# Workflows

An application developed for a backend test.

## Service Stack

For this test, the following technologies used
 - Golang as API language
 - gorm as an ORM framework with Postgres as the database driver
 - Shell Script for automation
 - Docker used as the containerization engine 
 - Docker Swarm for orchestrating the application containers
 - RabbitMQ used for messaging.

The idea of using RabbitMQ was for simplicity of deployment (a single container with a web dashboard for viewing and the server), but in the course of development, some specifications were not reached, for example, the consumption of a single message from the respective queue without destroying it, since RabbitMQ does not persist by default as Kafka does.
 

## Deployment

To facilitate the deployment of the service stack, the Docker and Docker Swarm was used. The deployment of the stack is done in three steps.

1 - Check if the Docker Swarm is activated on your local node. Be sure that you have at least Docker 19.03.9 installed and Swarm Mode activated by running:
 
```bash
docker swarm init
```
If you already have the swarm initialized, an error message as the following will occur

``` 
Error response from daemon: This node is already part of a swarm.
Use "docker swarm leave" to leave this swarm and join another one.
```
This means that you are ready to deploy the service stack using Docker/Docker Swarm.

2 - Fulfill the `.env` file properly, if they don't already exist you can check the `.env.example` file & copy and paste the content.

3 - Run
```
./stack/stack.sh deploy
```
If wild errors appear, just run again. If all is good, you will be able to access the endpoints. For tests purpose, 
run `$ curl http://127.0.0.1:5000/v1/version` or in whatever HTTP Client of your choice. 



# Backend Test

Develop the workflow's REST API following the specification bellow and document it.

## Defining a workflow

|Name|Type|Description|
|-|-|-|
|UUID|UUID|workflow unique identifier|
|status|Enum(inserted, consumed)|workflow status|
|data|JSONB|workflow input|
|steps|Array|name of workflow steps

## Endpoints

|Verb|URL|Description|
|-|-|-|
|POST|/workflow|insert a workflow on database and on queue and respond request with the inserted workflow|
|PATCH|/workflow/{UUID}|update status from specific workflow|
|GET|/workflow|list all workflows|
|GET|/workflow/consume|consume a workflow from queue and generate a CSV file with workflow.Data|

## Technologies

- Go, C, C++, Python, Java or any other that you know
- PostgreSQL
- A message queue that you choose, but describe why you choose.

