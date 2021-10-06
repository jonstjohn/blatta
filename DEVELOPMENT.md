# Development Notes

## Monitor Hot Ranges

Start cockroach demo:

```bash
cockroach demo
```

Monitor ranges

```bash
go run main.go monitor hotRanges --url='http://127.0.0.1:8080' \
  --username=demo --password=XXXXXX \
  --pgurl='postgres://demo:demo26350@127.0.0.1:26257?sslmode=require&sslrootcert=/var/folders/rg/_ltnftys4wj6ltzm0y45mhsh0000gq/T/demo449306277/ca.crt'
```