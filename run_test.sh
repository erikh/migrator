#!bash

# temporary so we can get it into tinyci before build pipelines land.

set -e

apt-get update
apt-get install postgresql-all -y

mkdir /postgresql
chown postgres:postgres /postgresql

su - postgres -c '/usr/lib/postgresql/9.6/bin/initdb -D /postgresql'
su - postgres -c '/usr/lib/postgresql/9.6/bin/pg_ctl -D /postgresql start &>/dev/null'

while ! (echo "select 1" | su - postgres sh -c 'psql template1 &>/dev/null')
do
  sleep 1
done

export IN_CONTAINER=1 # used by tests to determine if it can drop databases

su - postgres -c 'createuser -s root'

go get -t github.com/erikh/migrator/... && go test -race -v github.com/erikh/migrator -check.v
