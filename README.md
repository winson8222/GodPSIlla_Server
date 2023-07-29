![Getting Started](godzilla.png)

<div style="border: 1px solid white; padding: 10px;">

IDL_management_page contains the user interface

JSONGenwithpostres contains the Code generator

pages-backend-test contains the backend server for interacting with the database

</div>

Docker Setup guide

1.  RUN "docker pull quay.io/coreos/etcd:v3.5.0" to pull etcd image from dockerhub
2.  RUN "docker network create --driver bridge etcd-net" create a etcd-net network
3.  RUN "docker-compose up" at root directory to start he docker container for etcd
4.  RUN "docker pull postgres:latest" to pull etcd image from postgres
5.  RUN "docker run --name some-postgres -e POSTGRES_PASSWORD=mysecretpassword -e POSTGRES_DB=mydatabase -p 5432:5432 -d postgres" to run the database on port 5432
6.  Start the containers on docker, preferrably through docker desktop

Frontend Setup guide

1. RUN "cd .\IDL_management_page\page\my-app\"
2. RUN "npm install" to ensure that the dependency is up to date
3. RUN "npx prisma migrate dev --name init --schema=prisma/schema.prisma" to apply the migration using the schema to the database
4. RUN "npx prisma generate --schema=prisma/schema.prisma" generate a prisma client using the schema

Backend Setup guide

1. RETURN TO ROOT DIRECTORY
2. RUN "cd .\pages-backend-test\pages-backend\"
3. RUN "npx prisma migrate dev --name init --schema=prisma/schema.prisma" to apply the migration using the schema to the database
4. RUN "npx prisma generate --schema=prisma/schema.prisma" generate a prisma client using the schema

Nginx Setup guide (MacOS)

1. RUN "brew install nginx" install nginx on system

Nginx Setup guide (Windows)

1. Download Nginx from official Nginx website
2. RUN "setx NGINX_PATH "C:\path\to\nginx\directory" /M" on terminal (IMPORTANT) so as to make the nginx directory availble as a system environmental variable to be located to start nginx

To Start The Application:

1. RUN "npm run dev" under \IDL_management_page\page\my-app\ to run the frontend UI
2. RUN "node index.js" under \pages-backend-test\pages-backend\ to run the backend server
3. Open a browser on localhost:3000

Assumptions:

1. Docker is installed
2. Node.js is installed
3. Nginx is installed
4. Prisma installed globally using "npm install -g prisma"
5. No process is running on ports 8888, 8889, 8890, 20000, 5432, 3000, 3333 (Will make improvement to the customisibility in the future)
