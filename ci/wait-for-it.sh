#!/bin/ash
# wait-for-it.sh
set -e

until nc -z db 3306; do
  >&2 echo "mysql is unavailable - sleeping"
  sleep 3
done
>&2 echo "mysql is up - executing command"

exec $@

sleep 10
