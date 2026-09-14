#!/usr/bin/env bash
set -e

cd "$(dirname "$0")"

# Start compose in the background
docker compose up --build -d

# Wait for Delve to be reachable on 2345
echo "Waiting for Delve on localhost:2345..."
until nc -z localhost 2345; do
  sleep 0.5
done
echo "Delve is up."
sleep 1