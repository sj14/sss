# NAME

sss - S3 client

# SYNOPSIS

sss

```
[--access-key]=[value]
[--bandwidth]=[value]
[--bucket]=[value]
[--config]=[value]
[--endpoint]=[value]
[--header]=[value]
[--help|-h]
[--insecure]
[--path-style]
[--profile]=[value]
[--read-only]
[--region]=[value]
[--secret-key]=[value]
[--sni]=[value]
[--verbosity]=[value]
[--version|-v]
```

**Usage**:

```
sss [GLOBAL OPTIONS] [command [COMMAND OPTIONS]] [ARGUMENTS...]
```

# GLOBAL OPTIONS

**--access-key**="": 

**--bandwidth**="": Limit the bandwith per second (e.g. '1 MiB'). If set, an initial burst of 128 KiB is added.

**--bucket**="": 

**--config**="": ~/.config/sss/config.toml

**--endpoint**="": 

**--header**="": format: 'key1:val1,key2:val2'

**--help, -h**: show help

**--insecure**: 

**--path-style**: 

**--profile**="":  (default: "default")

**--read-only**: 

**--region**="": 

**--secret-key**="": 

**--sni**="": 

**--verbosity**="":  (default: 1)

**--version, -v**: print the version


# COMMANDS

## profiles

Config Profiles

**--help, -h**: show help

### help, h

Shows a list of commands or help for one command

## buckets

Bucket List

**--help, -h**: show help

**--prefix**="": 

### help, h

Shows a list of commands or help for one command

## bucket

Bucket Head

**--help, -h**: show help

### help, h

Shows a list of commands or help for one command

## mb

Bucket Create

**--help, -h**: show help

**--object-lock**: 

### help, h

Shows a list of commands or help for one command

## rb

Bucket Remove

**--help, -h**: show help

### help, h

Shows a list of commands or help for one command

## size

Bucket Size

**--delimiter**="":  (default: "/")

**--help, -h**: show help

### help, h

Shows a list of commands or help for one command

## policy

Bucket Policy

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

## versioning

Bucket Versioning

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

## object-lock

Bucket Object Locking

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

Bucket Lifecycle

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

## cors

Bucket CORS

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

## tag

Bucket Tagging

**--help, -h**: show help

### get


**--help, -h**: show help

#### help, h

Shows a list of commands or help for one command

### help, h

Shows a list of commands or help for one command

## multiparts

Multipart Uploads

**--help, -h**: show help

### ls


**--delimiter**="":  (default: "/")

**--help, -h**: show help

**--prefix**="": 

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

Multipart Parts

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

Object List

**--delimiter**="":  (default: "/")

**--help, -h**: show help

### help, h

Shows a list of commands or help for one command

## head

Object Head

**--help, -h**: show help

### help, h

Shows a list of commands or help for one command

## get

Object Download

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

## put

Object Upload

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

Object Remove

**--concurrency**="":  (default: 5)

**--delimiter**="":  (default: "/")

**--dry-run**: 

**--force**: 

**--help, -h**: show help

### help, h

Shows a list of commands or help for one command

## cp

Object Server Side Copy

**--dst-bucket**="": Destinaton bucket

**--dst-key**="": Destination key. When empty, the src-key will be used

**--help, -h**: show help

**--src-bucket**="": Source bucket

**--src-key**="": Source key

**--sse-c-algorithm**="":  (default: "AES256")

**--sse-c-key**="": 32 bytes key

### help, h

Shows a list of commands or help for one command

## versions

Object Versions

**--delimiter**="":  (default: "/")

**--help, -h**: show help

### help, h

Shows a list of commands or help for one command

## acl

Object ACL

**--help, -h**: show help

### get


**--help, -h**: show help

**--version-id**="": 

#### help, h

Shows a list of commands or help for one command

### help, h

Shows a list of commands or help for one command

## presign

Object pre-signed URL

**--expires-in**="":  (default: 0s)

**--help, -h**: show help

### get

Presigned URL for a GET request

**--help, -h**: show help

#### help, h

Shows a list of commands or help for one command

### put

Presigned URL for a PUT request

**--help, -h**: show help

#### help, h

Shows a list of commands or help for one command

### help, h

Shows a list of commands or help for one command

## help, h

Shows a list of commands or help for one command

## completion

Output shell completion script for bash, zsh, fish, or Powershell

**--help, -h**: show help

### help, h

Shows a list of commands or help for one command
