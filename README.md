# Shopify: inventory management system

## Instructions to run:

Run command: `docker-compose up -d`

Minimal frontend repository here: https://github.com/nuzurie/shopifyFE

It's possible frontend won't work since it's an image on docker hub and might be stale. In that case run the frontend
react app as in step 5. Or build the docker image and run it there.

By default, backend runs on port 8080, and frontend 3000. Change it in docker-compose file if needed.
Postgres db uses docker network so no need to change or expose anything.

If you wish to run them separately since FE is just the image:
1. `docker network create my-network`
2. `docker run -e POSTGRES_USER=docker -e POSTGRES_PASSWORD=docker -e POSTGRES_DB=shopify --network my-network -d -v /var/lib/postgresql/data --name postgres library/postgres`
3. `docker build -t backend .`
4. `docker run -e DATABASE_URL=postgres://docker:docker@postgres/shopify --network my-network -p 8080:8080 backend (this won't have sample entries added as in docker compose)`
5. Go to FE repository and run: `npm install`, and then `npm start`

Unfortunately, I'm not fast enough with Frontend tech to showcase all the backend edge-cases and design choices.
**Remove the default values to see placeholders for instructions.**
**Clicking on _Filter_ button will fetch new results in both item catalogue and inventory.**

## Architecure and Project Design

I used domain-driven design. Fun exercise: find all 5 SOLID principles in the project!

