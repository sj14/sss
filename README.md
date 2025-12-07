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
NAME:
   sss - S3 client

USAGE:
   sss [global options] [command [command options]]

COMMANDS:
   profiles    Config Profiles
   help, h     Shows a list of commands or help for one command
   completion  Output shell completion script for bash, zsh, fish, or Powershell

   bucket management:
     buckets      Bucket List
     bucket       Bucket Head
     mb           Bucket Create
     rb           Bucket Remove
     size         Bucket Size
     policy       Bucket Policy
     versioning   Bucket Versioning
     object-lock  Bucket Object Locking
     lifecycle    Bucket Lifecycle
     cors         Bucket CORS
     tag          Bucket Tagging

   multipart management:
     multiparts  Multipart Uploads
     parts       Multipart Parts

   object management:
     ls        Object List
     head      Object Head
     get       Object Download
     put       Object Upload
     rm        Object Remove
     cp        Object Server Side Copy
     versions  Object Versions
     acl       Object ACL
     presign   Object pre-signed URL

GLOBAL OPTIONS:
   --config string                      ~/.config/sss/config.toml [$SSS_CONFIG]
   --profile string                     (default: "default") [$SSS_PROFILE]
   --access-key string                   [$SSS_ACCESS_KEY]
   --secret-key string                   [$SSS_SECRET_KEY]
   --endpoint string                     [$SSS_ENDPOINT]
   --region string                       [$SSS_REGION]
   --path-style                          [$SSS_PATH_STYLE]
   --insecure                            [$SSS_INSECURE]
   --bucket string                       [$SSS_BUCKET]
   --read-only                           [$SSS_READ_ONLY]
   --bandwidth string                   Limit bandwith per second, e.g. '1 MiB' (always adds 128 KiB burst) [$SSS_BANDWIDTH]
   --sni string                          [$SSS_SNI]
   --header string [ --header string ]  format: 'key1:val1,key2:val2'
   --verbosity uint                     (default: 1) [$SSS_VERBOSITY]
   --help, -h                           show help
   --version, -v                        print the version
```

### Shell completion

Follow the instructions from `sss completion --help`.

### Examples

#### List objects

##### List bucket root

```
➜ sss --bucket <BUCKET> ls
                      PREFIX  test/
2025-11-22 11:11:05  100 MiB  100MB.bin
```

##### List directory/prefix

```
➜ sss --bucket <BUCKET> ls test/
2025-11-22 14:19:58  1.0 MiB  1MB.bin
2025-11-22 14:20:00  2.0 MiB  2MB.bin
```

#### Download

##### Download a single object

```
➜ test sss --bucket <BUCKET> get 100MB.bin
100 MiB in 11s | 9.0 MiB/s | 100MB.bin
```

##### Download directory/prefix

Only works when the end of the prefix matches the delimiter (by default: `/`).

```
➜ test sss --bucket <BUCKET> get test/
1.0 MiB in 0s | 2.5 MiB/s | test/1MB.bin
2.0 MiB in 0s | 5.0 MiB/s | test/2MB.bin
```

#### Upload

##### Upload a single object:

```
➜ sss --bucket <BUCKET> put 1MB.bin
1.0 MiB in 1s | 808 KiB/s | 1MB.bin                                             
```

##### Upload a directory:

```
➜ sss --bucket <BUCKET> put test/
1.0 MiB in 1s | 904 KiB/s | 1MB.bin                                             
2.0 MiB in 2s | 1.2 MiB/s | 2MB.bin                                             
```

Notice that the uploaded objects resist in the root of the bucket.

##### Upload to a specific path/prefix

```
➜ sss --bucket <BUCKET> put test/ test2/
1.0 MiB in 1s | 883 KiB/s | test2/1MB.bin
2.0 MiB in 2s | 1.0 MiB/s | test2/2MB.bin
```

#### Delete

##### Delete a single object

```
➜ sss --bucket <BUCKET> rm 100MB.bin
deleting 100MB.bin (100 MiB)
```

##### Delete a directory/prefix

Only works when the end of the prefix matches the delimiter (by default: `/`).

```
➜ sss --bucket <BUCKET> rm test/
deleting test/2MB.bin (2.0 MiB)
deleting test/1MB.bin (1.0 MiB)
```
