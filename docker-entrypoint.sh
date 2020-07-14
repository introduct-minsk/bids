#!/usr/bin/env bash
set -eo pipefail

DB_NAME=bids
PG_PASSWORD=postgres

if [[ -z "$PG_HOST" ]]
then
  PG_HOST=postgres
fi

if [[  -z "$PG_PORT" ]]
then
  PG_PORT=5432
fi

if [[ -z "$DB_NAME" ]]
then
  DB_NAME=bids
fi

if [[ -z "$PG_USER" ]]
then
  PG_USER=postgres
fi

if [[ -z "$PG_PASSWORD" ]]
then
  AUTH_STRING=${PG_USER}
else
  AUTH_STRING=${PG_USER}:${PG_PASSWORD}
fi

# ----------------------------------------------------------------------------------------
# Run migrations -------------------------------------------------------------------------
# ----------------------------------------------------------------------------------------
echo "Running migrations..."
n=0
timeout=15 #timeout seconds

until [[ ${n} -ge ${timeout} ]]
do
   ./app/go-migrate -database "postgresql://$AUTH_STRING@$PG_HOST:$PG_PORT/$DB_NAME?sslmode=disable" -path ./app/migrations up && break
    n=$[$n+1]
    if [[ ${n} -eq ${timeout} ]]
    then
        echo "Timeout reached. Exiting..."
        exit 1
    else
        sleep 1
        echo "Retrying..."
    fi
done

echo "Done."

exec "$@"
