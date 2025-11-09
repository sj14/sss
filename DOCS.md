# NAME

sss - S3 client

# SYNOPSIS

sss

```
[--access-key]=[value]
[--bucket]=[value]
[--endpoint]=[value]
[--help|-h]
[--insecure]
[--path-style]
[--profile]=[value]
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

**--endpoint**="": 

**--help, -h**: show help

**--insecure**: 

**--path-style**: 

**--profile**="":  (default: "default")

**--region**="": 

**--secret-key**="": 

**--verbosity**="":  (default: 1)


# COMMANDS

## buckets


**--help, -h**: show help

**--prefix**="": 

### help, h

Shows a list of commands or help for one command

## bucket


**--help, -h**: show help

### help, h

Shows a list of commands or help for one command

## mb


**--help, -h**: show help

### help, h

Shows a list of commands or help for one command

## rb


**--help, -h**: show help

### help, h

Shows a list of commands or help for one command

## multiparts


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


**--delimiter**="":  (default: "/")

**--help, -h**: show help

### help, h

Shows a list of commands or help for one command

## cp


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


**--concurrency**="":  (default: 5)

**--delimiter**="":  (default: "/")

**--force**: 

**--help, -h**: show help

### help, h

Shows a list of commands or help for one command

## get


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


**--help, -h**: show help

### help, h

Shows a list of commands or help for one command

## presign


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


**--help, -h**: show help

### get


**--help, -h**: show help

#### help, h

Shows a list of commands or help for one command

### help, h

Shows a list of commands or help for one command

## acl


**--help, -h**: show help

### get


**--help, -h**: show help

**--version-id**="": 

#### help, h

Shows a list of commands or help for one command

### help, h

Shows a list of commands or help for one command

## versions


**--delimiter**="":  (default: "/")

**--help, -h**: show help

### help, h

Shows a list of commands or help for one command

## help, h

Shows a list of commands or help for one command
