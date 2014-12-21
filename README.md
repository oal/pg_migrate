# pg_migrate

A very simple database migration tool for PostgreSQL. No tests yet - use at your own risk!

### How do I use it?
Create a directory to hold your migrations. Migrations must be named `[number]-[description].sql`, where number starts from 1. An example could be `1-initial.sql` and `2-add-post-table.sql` and so on.

`pg_migrate` will create a new table in your database called `migrations` to keep track of which migrations have already been applied. If you ran `pg_migrate` with the two migrations above, then later added `3-alter-post-table.sql`, `pg_migrate` would only apply this migration the next time it.

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
  --history            Show migration history
  --verbose            Verbose output
```