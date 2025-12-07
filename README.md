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
   acl          Handle Object ACL
   versions     List Object Versions
   help, h      Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --endpoint string
   --insecure
   --region string
   --path-style
   --profile string
   --bucket string
   --secret-key string
   --access-key string
   --verbosity uint     (default: 1)
   --help, -h           show help
```

## Configuraton

`sss` uses the AWS package and the same configuration files and environment variables (e.g. `AWS_ACCESS_KEY_ID` and `AWS_SECRET_ACCESS_KEY`) as the AWS CLI.

`~/.aws/config`:

```ini
[default]
endpoint_url = https://example.com
region = earth
```

`~/.aws/credentials`:

```ini
[default]
aws_access_key_id = <CHANGE_ME>
aws_secret_access_key = <CHANGE_ME>
```
