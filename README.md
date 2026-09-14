# Stop and remove containers + networks created by this compose project
docker compose down

# Also remove the images built for this project
docker compose down --rmi local

docker image prune -f

docker compose up --build