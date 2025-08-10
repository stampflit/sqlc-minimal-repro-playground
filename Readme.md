Corresponding playground: https://play.sqlc.dev/p/bd2b6da368591bfd7277bf834deda98e2820689469b7d8e6cf3adc8b29ddfb49


Start a db

```sh
docker run --rm -e POSTGRES_PASSWORD=password -e POSTGRES_USER=user -e POSTGRES_DB=testdb -p 5433:5432 postgres
```

Run the migration

```sh
PGPASSWORD='password' psql -h localhost -p 5433 -U user -d testdb -f ./schema.sql
```

See

- https://dishtaxonomy.com/
- https://saladtheory.github.io/
- https://archive.is/cDdpL
