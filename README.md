# pg_migrate

A very simple database migration tool for PostgreSQL. No tests yet - use at your own risk!

```
$ pg_migrate --help
usage: pg_migrate [<flags>]

Flags:
  --help               Show help.
  --dir="/Users/olav/Sources/Go/src/github.com/oal/pg_migrate"
                       Directory where migration files are located. Current working directory will be
                       used by default
  --host="localhost"   Server address and port
  --port=5432          Server port
  --db="postgres"      Database name
  --user="postgres"    User
  --password=PASSWORD  Password
  --sslmode="disable"
  --verbose            Verbose output
```