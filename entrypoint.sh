#!/bin/sh

su - postgres -c '/usr/lib/postgresql/9.6/bin/pg_ctl -D /postgresql start &>/dev/null'

while ! (echo "select 1" | su - postgres sh -c 'psql template1 &>/dev/null')
do
  sleep 1
done

su - postgres -c 'createuser migrator'
su - postgres -c 'createdb migrator -O migrator'

exec "$@"
