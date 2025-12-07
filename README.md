# sss

Yet another S3 client.

## Documentation

[DOCS.md](DOCS.md) contains the generate documentation.

## Usage

For shell completion follow the instructions from `sss completion --help`.

```
NAME:
   sss object-lock get

USAGE:
   sss object-lock get [options]

OPTIONS:
   --help, -h  show help

GLOBAL OPTIONS:
   --config string                      ~/.config/sss/config.yaml [$SSS_CONFIG]
   --endpoint string                     [$SSS_ENDPOINT]
   --insecure                            [$SSS_INSECURE]
   --read-only                           [$SSS_READ_ONLY]
   --region string                       [$SSS_REGION]
   --path-style                          [$SSS_PATH_STYLE]
   --profile string                     (default: "default") [$SSS_PROFILE]
   --bucket string
   --secret-key string                   [$SSS_SECRET_KEY]
   --access-key string                   [$SSS_ACCESS_KEY]
   --sni string                          [$SSS_SNI]
   --header string [ --header string ]  format: 'key:val'
   --verbosity uint                     (default: 1) [$SSS_VERBOSITY]
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
  # path_style: true
  # read_only: true
```
