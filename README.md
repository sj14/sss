# sss

Yet another S3 client.

## Installation

### Binaries

Binaries are available for all major platforms. See the [releases](https://github.com/sj14/sss/releases) page.

### Container

```bash
docker pull ghcr.io/sj14/sss
```

### Homebrew

Using the [Homebrew](https://brew.sh/) package manager for macOS:

``` text
brew install sj14/tap/sss
```

### Go

It's also possible to install via `go install`:

```console
go install github.com/sj14/sss@latest
```

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
      --config=STRING        Path to the config file (default:
                             ~/.config/sss/config.toml) ($SSS_CONFIG).
      --profile="default"    Profile to use ($SSS_PROFILE).
      --endpoint=STRING      S3 endpoint URL ($SSS_ENDPOINT).
      --region=STRING        S3 region ($SSS_REGION).
      --path-style           Use path style S3 requests ($SSS_PATH_STYLE).
      --access-key=STRING    S3 access key ($SSS_ACCESS_KEY).
      --secret-key=STRING    S3 secret key ($SSS_SECRET_KEY).
      --insecure             Skip TLS verification ($SSS_INSECURE).
      --read-only            Only allow safe HTTP methods ($SSS_READ_ONLY).
      --bandwidth=STRING     Limit bandwith per second, e.g. '1 MiB' (always 64
                             KiB burst) ($SSS_BANDWIDTH).
      --header=HEADER,...    Add additional HTTP headers (format:
                             'key1:val1,key2:val2') ($SSS_HEADER).
      --sni=STRING           TLS Server Name Indication ($SSS_SNI).

Commands:
  profiles [flags]
    List availale profiles.

  version [flags]
    Show version information.

Bucket Commands
  buckets [flags]
    List all buckets.

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

Multipart Commands
  bucket <bucket> multiparts rm <object> <upload-id>

  bucket <bucket> multiparts ls [<prefix>] [flags]

  bucket <bucket> multiparts parts ls <object> <upload-id> [flags]

Object Commands
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

### Delimiter

The forward slash (`/`) is the only supported delimiter.

### Examples

Instead of using `sss bucket <BUCKET>` you can also use the shorter variant `sss b <BUCKET>`.

#### List Buckets

```
➜ sss buckets
2025-11-07 12:54:43 <bucket-a>
2024-12-19 09:16:14 <bucket-b>
2025-11-20 14:49:38 <bucket-c>
```

Or use `sss ls` as the shorter alternative.

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

##### Download a single object with different name

```
➜ test sss bucket <BUCKET> get 1MB.bin yolo
1.0 MiB in  0s |  19 MiB/s | yolo
```

##### Download a single object into a directory

```
➜ sss bucket test get 1MB.bin test/
1.0 MiB in  0s |  19 MiB/s | test/1MB.bin
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

##### Upload a single object with different name:

```
➜ sss bucket simon-crasher put ~/Downloads/1MB.bin yolo
1.0 MiB in  1s | 1.6 MiB/s | yolo
```

##### Upload a single object into a directory/prefix:

```
➜ sss bucket simon-crasher put ~/Downloads/1MB.bin yolo/
1.0 MiB in  1s | 1.6 MiB/s | yolo/1MB.bin
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
