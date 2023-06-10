# go-safe

Backup your data safely, and easily.

## CLI Usage 

# Backup

You can run the backup utility in foreground : 

```go
go build -o go-safe cmd/docker/*
./go-safe --access-id=**** --access-key=**** --backup-dir=backup --bucket-name=**** --endpoint=**** --region=**** --interval=10 --rsa-pubkey-location=public.pem 
```