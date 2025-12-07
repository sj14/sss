# sss

Yet another S3 client.

## Documentation

[DOCS.md](DOCS.md) contains the generate documentation.

## Usage

For shell completion follow the instructions from `sss complete`. For example, add `source <(sss completion zsh)` to your `~/.zshrc`.


```
NAME:
   sss - S3 client

USAGE:
   sss [global options] [command [command options]]

COMMANDS:
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

GLOBAL OPTIONS:
   --config string      ~/.config/sss/config.yaml [$SSS_CONFIG]
   --endpoint string     [$SSS_ENDPOINT]
   --insecure            [$SSS_INSECURE]
   --read-only           [$SSS_READ_ONLY]
   --region string       [$SSS_REGION]
   --path-style          [$SSS_PATH_STYLE]
   --profile string     (default: "default") [$SSS_PROFILE]
   --bucket string
   --secret-key string   [$SSS_SECRET_KEY]
   --access-key string   [$SSS_ACCESS_KEY]
   --verbosity uint     (default: 1) [$SSS_VERBOSITY]
   --help, -h           show help
```

## Configuraton

`~/.config/sss/config.yaml`:

```yaml
profiles:
  default:
    endpoint: https://example.com
    region: earth
    access_key: <CHANGE_ME>
    secret_key: <CHANGE_ME>
```
