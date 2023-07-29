## Postresg setup:

- docker run --name some-postgres -e POSTGRES_PASSWORD=mysecretpassword -p 5432:5432 -d postgres
- npx prisma migrate dev --name init --schema=pages-backend/prisma/schema.prisma
- npx prisma generate --schema=pages-backend/prisma/schema.prisma

## Backend Setup:

- cd pages-backend
- node index.js

## Changes:

- Included IDL filename in Schema
