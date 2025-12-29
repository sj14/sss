# sss

`sss` is yet another S3 client. I respect all the existing clients and `sss` certainly has its own shortcomings, but some reasons for starting the project:

- Having a single deployable binary.
- A CLI and config which can handle non-AWS S3 use-cases well ([compare](https://github.com/aws/aws-cli/issues/4215)).
- Get somewhat meaningful errors instead of `'NoneType' is not iterable`.
- Having a read-only mode (allowing only [safe](https://httpwg.org/specs/rfc9110.html#safe.methods) HTTP methods).
- Having the S3 feature set I need.
- Be able to limit the bandwidth.
- Ergonomic usage, e.g. `<bucket> get <object> --sse-c-key=<key>` ([compare](https://docs.min.io/enterprise/aistor-object-store/reference/cli/mc-put/#--enc-c)).
- Hard to misuse (no local file management like `mc cp /mydata/ alias-typo/mybucket/mydata`).
- Writing the bucket name before the operation to switch the operation quickly, e.g. `<bucket> policy get` and `<bucket> policy rm`.


## Installation

### Binaries

Binaries are available for all major platforms. See the [releases](https://github.com/sj14/sss/releases) page.

### Homebrew

Using the [Homebrew](https://brew.sh/) package manager for macOS:

``` text
brew install sj14/tap/sss
```

### Container

```bash
docker pull ghcr.io/sj14/sss
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
network    = "tcp6"
bandwidth  = "128 MiB"
```

## Usage

```
Usage: sss <command> [flags]

Generic Commands
  config (c) show (s)        Get config.
  config (c) profiles (p)    List availale profiles.
  version                    Show version information.

Bucket Commands
  buckets (ls)                                List all buckets.
  bucket (b) <bucket> mb                      Make/create bucket.
  bucket (b) <bucket> hb                      Head bucket lists bucket information.
  bucket (b) <bucket> rb                      Remove/delete bucket.
  bucket (b) <bucket> policy get              Get lifecycle policy.
  bucket (b) <bucket> policy put              Put lifecycle policy.
  bucket (b) <bucket> policy rm               Delete lifecycle policy.
  bucket (b) <bucket> cors get                Get CORS policy.
  bucket (b) <bucket> cors put                Put CORS policy.
  bucket (b) <bucket> cors rm                 Delete CORS policy.
  bucket (b) <bucket> tag get                 Get bucket tag.
  bucket (b) <bucket> lifecycle (lc) get      Get lifecycle policy.
  bucket (b) <bucket> lifecycle (lc) put      Put lifecycle policy.
  bucket (b) <bucket> lifecycle (lc) rm       Delte lifecycle policy.
  bucket (b) <bucket> versioning get          Get bucket versioning config.
  bucket (b) <bucket> versioning put          Put bucket versioning config.
  bucket (b) <bucket> cleanup                 Remove all objects versions and multiparts from the bucket.
  bucket (b) <bucket> object-lock (ol) get    Get object-lock config.
  bucket (b) <bucket> object-lock (ol) put    Put object-lock config.
  bucket (b) <bucket> size                    Calculate bucket size (resource heavy!)

Object Commands
  bucket (b) <bucket> ls             List objects.
  bucket (b) <bucket> cp             Server-side copy.
  bucket (b) <bucket> put            Upload object(s).
  bucket (b) <bucket> put-rand       Upload random object(s).
  bucket (b) <bucket> rm             Remove object.
  bucket (b) <bucket> get            Download object(s). Requires HeadObject permission.
  bucket (b) <bucket> head           Head Object Liss object information.
  bucket (b) <bucket> versions       List object versions
  bucket (b) <bucket> presign get    Create pre-signed URL for GET request.
  bucket (b) <bucket> presign put    Create pre-signed URL for PUT request.
  bucket (b) <bucket> acl get        Get object ACL.

Multipart Commands
  bucket (b) <bucket> multipart (mp) rm          Delete multipart upload.
  bucket (b) <bucket> multipart (mp) ls          List multipart uploads.
  bucket (b) <bucket> multipart (mp) create      Create multipart upload.
  bucket (b) <bucket> multipart (mp) parts ls    List parts.

Flags:
  -h, --help                    Show context-sensitive help.
  -c, --config=STRING           Path to the config file (default: ~/.config/sss/config.toml) ($SSS_CONFIG).
  -p, --profile="default"       Profile to use ($SSS_PROFILE).
  -v, --verbosity=1             Output verbosity (0=disable; 1=default; 8=header; 9=body) ($SSS_VERBOSITY).
      --endpoint=STRING         S3 endpoint URL ($SSS_ENDPOINT).
      --region=STRING           S3 region ($SSS_REGION).
      --path-style              Use path style S3 requests ($SSS_PATH_STYLE).
      --access-key=STRING       S3 access key ($SSS_ACCESS_KEY).
      --secret-key=STRING       S3 secret key ($SSS_SECRET_KEY).
      --insecure                Skip TLS verification ($SSS_INSECURE).
      --read-only               Only allow safe HTTP methods (HEAD, GET, OPTIONS) ($SSS_READ_ONLY).
      --network="tcp"           Force IPv4/6 with 'tcp4' or 'tcp6' ($SSS_NETWORK).
      --bandwidth=STRING        Limit bandwith per second, e.g. '1 MiB' (always 64 KiB burst) ($SSS_BANDWIDTH).
      --header=KEY=VALUE;...    Set HTTP headers (format: 'key1=val1;key2=val2') ($SSS_HEADER).
      --param=KEY=VALUE;...     Set URL parameters (format: 'key1=val1;key2=val2') ($SSS_PARAM).
      --sni=STRING              TLS Server Name Indication ($SSS_SNI).

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
➜ sss bucket <bucket> put ~/Downloads/1MB.bin yolo
1.0 MiB in  1s | 1.6 MiB/s | yolo
```

##### Upload a single object into a directory/prefix:

Only works when the end of the prefix is `/`.

```
➜ sss bucket <bucket> put ~/Downloads/1MB.bin yolo/
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
