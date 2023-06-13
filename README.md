# go-safe

Backup your data safely, and easily.

## CLI Usage

### Backup

#### Flags

The flags for the CLI are the followings :

- `--s3.access-id`: S3 access ID
- `--s3.access-key`: S3 access key
- `--s3.bucket-name`: S3 bucket name
- `--s3.endpoint`: S3 endpoint
- `--s3.region`: S3 region
- `--s3.dir`: S3 directory (will store under a directory in S3)
- `--s3.storage-class`: S3 storage class
- `--backup.dir`: Backup directory
- `--interval`: Backup interval in seconds

And one of :

- `--aes.key-location`: AES key location
- `--ecies.public-key-location`: ECIES public key location

## Docker Usage

### Environment

The following environment variables are available

- S3 Access ID: GS_S3_ACCESS_ID
- S3 Access Key: GS_S3_ACCESS_KEY
- S3 Bucket Name: GS_S3_BUCKET_NAME
- S3 Endpoint: GS_S3_ENDPOINT
- S3 Region: GS_S3_REGION
- S3 Directory: GS_S3_DIR
- S3 Storage Class: GS_S3_STORAGE_CLASS
- AES Key Location: GS_AES_KEY_LOCATION
- ECIES Public Key Location: GS_ECIES_PUBLIC_KEY_LOCATION
- Backup Directory: GS_BACKUP_DIR
- Backup Interval: GS_INTERVAL

To use the backup tool properly, you must mount the `GS_BACKUP_DIR` and the encryption key of your liking.

The encryption key must be user-readable only. (`chmod 0400 key`)

### Export config

You can export your config if you need to use the retriever binary. To do, you can use the flag `--export` on the `go-safe` binary in the docker image.

## ECIES

To generate a compatible ECIES keypair, you can use the ecies-keygen utility provided in the different releases.

To use it, simply run `./ecies-keygen`.
