# sss

Yet another S3 client.

## Documentation

[DOCS.md](DOCS.md) contains the generate documentation.

## Usage

For shell completion follow the instructions from `sss completion --help`.

```
NAME:
   sss - S3 client

USAGE:
   sss [global options] [command [command options]]

COMMANDS:
   profiles     List config profiles
   buckets      List Buckets
   bucket       Head Bucket
   mb           Make Bucket
   rb           Remove Bucket
   multiparts   Handle Multipart Uploads
   parts        Parts from Multipart Uploads
   ls           List Objects
   cp           Server Side Object Copy
   put          Upload Object
   rm           Remove Object
   get          Download Object
   head         Head Object
   presign      Create pre-signed URL
   policy       Handle Bucket Policy
   cors         Handle Bucket CORS
   object-lock  Handle Bucket Object Locking
   lifecycle    Handle Bucket Lifecycle
   versioning   Handle Bucket Versioning
   size         Calculate the bucket size
   acl          Handle Object ACL
   versions     List Object Versions
   help, h      Shows a list of commands or help for one command
   completion   Output shell completion script for bash, zsh, fish, or Powershell

GLOBAL OPTIONS:
   --config string                      ~/.config/sss/config.toml [$SSS_CONFIG]
   --endpoint string                     [$SSS_ENDPOINT]
   --insecure                            [$SSS_INSECURE]
   --read-only                           [$SSS_READ_ONLY]
   --region string                       [$SSS_REGION]
   --path-style                          [$SSS_PATH_STYLE]
   --profile string                     (default: "default") [$SSS_PROFILE]
   --bucket string                       [$SSS_BUCKET]
   --secret-key string                   [$SSS_SECRET_KEY]
   --access-key string                   [$SSS_ACCESS_KEY]
   --sni string                          [$SSS_SNI]
   --header string [ --header string ]  format: 'key1:val1,key2:val2'
   --verbosity uint                     (default: 1) [$SSS_VERBOSITY]
   --help, -h                           show help
   --version, -v                        print the version
```

## Configuraton

`~/.config/sss/config.toml`:

```toml
[profiles.default]
endpoint = https://earth.example.com
region = earth
access_key = <CHANGE_ME>
secret_key = <CHANGE_ME>

[profiles.mars]
endpoint = https://mars.example.com
region = mars
access_key = <CHANGE_ME>
secret_key = <CHANGE_ME>
path_style = true
insecure   = true
read_only  = true
```
