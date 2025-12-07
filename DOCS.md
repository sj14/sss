# NAME

sss - S3 client

# SYNOPSIS

sss

```
[--access-key]=[value]
[--bucket]=[value]
[--config]=[value]
[--endpoint]=[value]
[--help|-h]
[--insecure]
[--path-style]
[--profile]=[value]
[--read-only]
[--region]=[value]
[--secret-key]=[value]
[--verbosity]=[value]
```

**Usage**:

```
sss [GLOBAL OPTIONS] [command [COMMAND OPTIONS]] [ARGUMENTS...]
```

# GLOBAL OPTIONS

**--access-key**="": 

**--bucket**="": 

**--config**="":  (default: "~/.config/sss/config.yaml")

**--endpoint**="": 

**--help, -h**: show help

**--insecure**: 

**--path-style**: 

**--profile**="":  (default: "default")

**--read-only**: 

**--region**="": 

**--secret-key**="": 

**--verbosity**="":  (default: 1)


# COMMANDS

## buckets

List Buckets

**--help, -h**: show help

**--prefix**="": 

### help, h

Shows a list of commands or help for one command

## bucket

Head Bucket

**--help, -h**: show help

### help, h

Shows a list of commands or help for one command

## mb

Make Bucket

**--help, -h**: show help

### help, h

Shows a list of commands or help for one command

## rb

Remove Bucket

**--help, -h**: show help

### help, h

Shows a list of commands or help for one command

## multiparts

Handle Multipart Uploads

**--help, -h**: show help

### ls


**--help, -h**: show help

#### help, h

Shows a list of commands or help for one command

### rm


**--help, -h**: show help

**--key**="": 

**--upload-id**="": 

#### help, h

Shows a list of commands or help for one command

### help, h

Shows a list of commands or help for one command

## parts

Parts from Multipart Uploads

**--help, -h**: show help

**--key**="": 

**--upload-id**="": 

### ls


**--help, -h**: show help

#### help, h

Shows a list of commands or help for one command

### help, h

Shows a list of commands or help for one command

## ls

List Objects

**--delimiter**="":  (default: "/")

**--help, -h**: show help

### help, h

Shows a list of commands or help for one command

## cp

Server Side Object Copy

**--dst-bucket**="": 

**--dst-key**="": when empty, the src-key will be used

**--help, -h**: show help

**--src-bucket**="": 

**--src-key**="": 

**--sse-c-algorithm**="":  (default: "AES256")

**--sse-c-key**="": 32 bytes key

### help, h

Shows a list of commands or help for one command

## put

Upload Object

**--acl**="": e.g. 'public-read'

**--concurrency**="":  (default: 5)

**--help, -h**: show help

**--leave-parts-on-error**="":  (default: 0)

**--max-parts**="":  (default: 0)

**--part-size**="":  (default: 0)

**--sse-c-algorithm**="":  (default: "AES256")

**--sse-c-key**="": 32 bytes key

**--target**="": target key for single file or prefix multiple files

### help, h

Shows a list of commands or help for one command

## rm

Remove Object

**--concurrency**="":  (default: 5)

**--delimiter**="":  (default: "/")

**--force**: 

**--help, -h**: show help

### help, h

Shows a list of commands or help for one command

## get

Download Object

**--concurrency**="":  (default: 5)

**--help, -h**: show help

**--if-match**="": 

**--if-modified-since**="": 

**--if-none-match**="": 

**--if-unmodified-since**="": 

**--part-number**="":  (default: 0)

**--part-size**="":  (default: 0)

**--range**="": bytes=BeginByte-EndByte, e.g. 'bytes=0-500' to get the first 501 bytes

**--sse-c-algorithm**="":  (default: "AES256")

**--sse-c-key**="": 32 bytes key

**--version-id**="": 

### help, h

Shows a list of commands or help for one command

## head

Head Object

**--help, -h**: show help

### help, h

Shows a list of commands or help for one command

## presign

Create pre-signed URL

**--expires-in**="":  (default: 0s)

**--help, -h**: show help

### get


**--help, -h**: show help

#### help, h

Shows a list of commands or help for one command

### put


**--help, -h**: show help

#### help, h

Shows a list of commands or help for one command

### help, h

Shows a list of commands or help for one command

## policy

Handle Bucket Policy

**--help, -h**: show help

### get


**--help, -h**: show help

#### help, h

Shows a list of commands or help for one command

### put


**--help, -h**: show help

#### help, h

Shows a list of commands or help for one command

### help, h

Shows a list of commands or help for one command

## cors

Handle Bucket CORS

**--help, -h**: show help

### get


**--help, -h**: show help

#### help, h

Shows a list of commands or help for one command

### put


**--help, -h**: show help

#### help, h

Shows a list of commands or help for one command

### rm


**--help, -h**: show help

#### help, h

Shows a list of commands or help for one command

### help, h

Shows a list of commands or help for one command

## object-lock

Handle Bucket Object Locking

**--help, -h**: show help

### get


**--help, -h**: show help

#### help, h

Shows a list of commands or help for one command

### put


**--help, -h**: show help

#### help, h

Shows a list of commands or help for one command

### help, h

Shows a list of commands or help for one command

## lifecycle

Handle Bucket Lifecycle

**--help, -h**: show help

### get


**--help, -h**: show help

#### help, h

Shows a list of commands or help for one command

### put


**--help, -h**: show help

#### help, h

Shows a list of commands or help for one command

### rm


**--help, -h**: show help

#### help, h

Shows a list of commands or help for one command

### help, h

Shows a list of commands or help for one command

## versioning

Handle Bucket Versioning

**--help, -h**: show help

### get


**--help, -h**: show help

#### help, h

Shows a list of commands or help for one command

### help, h

Shows a list of commands or help for one command

## acl

Handle Object ACL

**--help, -h**: show help

### get


**--help, -h**: show help

**--version-id**="": 

#### help, h

Shows a list of commands or help for one command

### help, h

Shows a list of commands or help for one command

## versions

List Object Versions

**--delimiter**="":  (default: "/")

**--help, -h**: show help

### help, h

Shows a list of commands or help for one command

## help, h

Shows a list of commands or help for one command
