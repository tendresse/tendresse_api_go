# Tendresse APP

```
docker run --name postgres-tendresse -e POSTGRES_PASSWORD=tendresse_password -e POSTGRES_USER=tendresse_user -e POSTGRES_DB=tendresse_app -p 127.0.0.1:5432:5432 -d postgres:10.5
```

---

```
# DATABASE_URL=postgres://$USERNAME:$PASSWORD@$HOST:$PORT/$DATABASE
```