from "golang:latest"

run "apt-get update -qq"
run "apt-get install postgresql-all -y"

run "mkdir /postgresql"
run "chown postgres:postgres /postgresql"

run "su - postgres -c '/usr/lib/postgresql/*/bin/initdb -D /postgresql'"

copy "entrypoint.sh", "/"
run "chmod 755 /entrypoint.sh"

entrypoint "/entrypoint.sh"
cmd "bash"
