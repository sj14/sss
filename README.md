# sss

Yet another S3 client.

## Documentation

[DOCS.md](DOCS.md) contains a generated documentation.

## Config File

`~/.config/sss/config.toml`:

```toml
[profiles.default]
endpoint = "https://earth.example.com"
region = "earth"
access_key = "<CHANGE_ME>"
secret_key = "<CHANGE_ME>"

[profiles.mars]
endpoint = "https://mars.example.com"
region = "mars"
access_key = "<CHANGE_ME>"
secret_key = "<CHANGE_ME>"
path_style = true
insecure   = true
read_only  = true
```

## Usage

```
Usage: sss <command> [flags]

Flags:
  -h, --help                 Show context-sensitive help.
      --config=STRING        ($SSS_CONFIG)
      --profile="default"    ($SSS_PROFILE)
      --endpoint=STRING      ($SSS_ENDPOINT)
      --region=STRING        ($SSS_REGION)
      --path-style           ($SSS_PATH_STYLE)
      --access-key=STRING    ($SSS_ACCESS_KEY)
      --secret-key=STRING    ($SSS_SECRET_KEY)
      --insecure             ($SSS_INSECURE)
      --read-only            ($SSS_READ_ONLY)
      --bandwidth=STRING     ($SSS_BANDWIDTH)
      --header=HEADER,...    ($SSS_HEADER)
      --sni=STRING           ($SSS_SNI)

Commands:
  profiles [flags]

  buckets [flags]

bucket
  bucket <bucket> mb [flags]

  bucket <bucket> hb

  bucket <bucket> rb [flags]

  bucket <bucket> policy get

  bucket <bucket> policy put <path>

  bucket <bucket> policy rm

  bucket <bucket> cors get

  bucket <bucket> cors put <path>

  bucket <bucket> cors rm

  bucket <bucket> tag get

  bucket <bucket> lifecycle get

  bucket <bucket> lifecycle put <path>

  bucket <bucket> lifecycle rm

  bucket <bucket> versioning get

  bucket <bucket> versioning put <path>

  bucket <bucket> object-lock get

  bucket <bucket> object-lock put <path>

  bucket <bucket> size [<path>]

multiparts
  bucket <bucket> multiparts rm <object> <upload-id>

  bucket <bucket> multiparts ls [<prefix>] [flags]

  bucket <bucket> multiparts parts ls <object> <upload-id> [flags]

object
  bucket <bucket> ls [<prefix>] [flags]

  bucket <bucket> cp <src-object> <dst-bucket> <dst-object>

  bucket <bucket> put <path> [<destination>] [flags]

  bucket <bucket> rm <object> [flags]

  bucket <bucket> get <object> [<destination>] [flags]

  bucket <bucket> head <object>

  bucket <bucket> presign get <object> [flags]

  bucket <bucket> presign put <object> [flags]

  bucket <bucket> acl get <object> [flags]

  bucket <bucket> versions [<path>] [flags]

Run "sss <command> --help" for more information on a command.
```

### Shell completion

Follow the instructions from `sss completion --help`.

### Delimiter

The forward slash (`/`) is the only supported delimiter.

### Examples

#### List objects

##### List bucket root

```
➜ sss bucket <BUCKET> ls
                      PREFIX  test/
2025-11-22 11:11:05  100 MiB  100MB.bin
```

##### List recursively

```
➜ sss bucket <BUCKET> ls -r
2025-11-22 14:19:58  1.0 MiB  test/1MB.bin
2025-11-22 14:20:00  2.0 MiB  test/2MB.bin
2025-11-22 11:11:05  100 MiB  100MB.bin
```

##### List directory/prefix

```
➜ sss bucket <BUCKET> ls test/
2025-11-22 14:19:58  1.0 MiB  1MB.bin
2025-11-22 14:20:00  2.0 MiB  2MB.bin
```

#### Download

##### Download a single object

```
➜ test sss bucket <BUCKET> get 100MB.bin
100 MiB in 11s | 9.0 MiB/s | 100MB.bin
```

##### Download directory/prefix

Only works when the end of the prefix is `/`.

```
➜ test sss bucket <BUCKET> get test/
1.0 MiB in 0s | 2.5 MiB/s | test/1MB.bin
2.0 MiB in 0s | 5.0 MiB/s | test/2MB.bin
```

#### Upload

##### Upload a single object:

```
➜ sss bucket <BUCKET> put 1MB.bin
1.0 MiB in 1s | 808 KiB/s | 1MB.bin                                             
```

##### Upload a directory:

```
➜ sss bucket <BUCKET> put test/
1.0 MiB in 1s | 904 KiB/s | test/1MB.bin                                             
2.0 MiB in 2s | 1.2 MiB/s | test/2MB.bin                                             
```

#### Delete

##### Delete a single object

```
➜ sss bucket <BUCKET> rm 100MB.bin
deleting 100MB.bin (100 MiB)
```

##### Delete a directory/prefix

Only works when the end of the prefix is `/`.

```
➜ sss bucket <BUCKET> rm test/
deleting test/2MB.bin (2.0 MiB)
deleting test/1MB.bin (1.0 MiB)
```
