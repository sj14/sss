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

## Shell completion

Follow the instructions from `sss completion --help`.

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
   --bandwidth string                   Limit the bandwith per second (e.g. '1 MiB'). If set, an initial burst of 128 KiB is added. [$SSS_BANDWIDTH]
   --sni string                          [$SSS_SNI]
   --header string [ --header string ]  format: 'key1:val1,key2:val2'
   --verbosity uint                     (default: 1) [$SSS_VERBOSITY]
   --help, -h                           show help
   --version, -v                        print the version
```
