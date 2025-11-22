## CLI interface - sss

S3 client.

Usage:

```bash
$ docs [GLOBAL FLAGS] [COMMAND] [COMMAND FLAGS] [ARGUMENTS...]
```

Global flags:

| Name               | Description                                                         | Type   | Default value | Environment variables |
|--------------------|---------------------------------------------------------------------|--------|:-------------:|:---------------------:|
| `--config="…"`     | ~/.config/sss/config.toml                                           | string |               |     `SSS_CONFIG`      |
| `--profile="…"`    |                                                                     | string |  `"default"`  |     `SSS_PROFILE`     |
| `--access-key="…"` |                                                                     | string |               |   `SSS_ACCESS_KEY`    |
| `--secret-key="…"` |                                                                     | string |               |   `SSS_SECRET_KEY`    |
| `--endpoint="…"`   |                                                                     | string |               |    `SSS_ENDPOINT`     |
| `--region="…"`     |                                                                     | string |               |     `SSS_REGION`      |
| `--path-style`     |                                                                     | bool   |    `false`    |   `SSS_PATH_STYLE`    |
| `--insecure`       |                                                                     | bool   |    `false`    |    `SSS_INSECURE`     |
| `--bucket="…"`     |                                                                     | string |               |     `SSS_BUCKET`      |
| `--read-only`      |                                                                     | bool   |    `false`    |    `SSS_READ_ONLY`    |
| `--bandwidth="…"`  | Limit bandwith per second, e.g. '1 MiB' (always adds 128 KiB burst) | string |               |    `SSS_BANDWIDTH`    |
| `--sni="…"`        |                                                                     | string |               |       `SSS_SNI`       |
| `--header="…"`     | format: 'key1:val1,key2:val2'                                       | string |               |        *none*         |
| `--verbosity="…"`  |                                                                     | uint   |      `1`      |    `SSS_VERBOSITY`    |
| `--help` (`-h`)    | show help                                                           | bool   |    `false`    |        *none*         |
| `--version` (`-v`) | print the version                                                   | bool   |    `false`    |        *none*         |

### `profiles` command

Config Profiles.

Usage:

```bash
$ docs [GLOBAL FLAGS] profiles [COMMAND FLAGS] [ARGUMENTS...]
```

The following flags are supported:

| Name            | Description | Type | Default value | Environment variables |
|-----------------|-------------|------|:-------------:|:---------------------:|
| `--help` (`-h`) | show help   | bool |    `false`    |        *none*         |

### `profiles help` subcommand (aliases: `h`)

Shows a list of commands or help for one command.

Usage:

```bash
$ docs [GLOBAL FLAGS] profiles help [command]
```

### `buckets` command

Bucket List.

Usage:

```bash
$ docs [GLOBAL FLAGS] buckets [COMMAND FLAGS] [ARGUMENTS...]
```

The following flags are supported:

| Name            | Description | Type   | Default value | Environment variables |
|-----------------|-------------|--------|:-------------:|:---------------------:|
| `--prefix="…"`  |             | string |               |        *none*         |
| `--help` (`-h`) | show help   | bool   |    `false`    |        *none*         |

### `buckets help` subcommand (aliases: `h`)

Shows a list of commands or help for one command.

Usage:

```bash
$ docs [GLOBAL FLAGS] buckets help [command]
```

### `bucket` command

Bucket Head.

Usage:

```bash
$ docs [GLOBAL FLAGS] bucket [COMMAND FLAGS] [ARGUMENTS...]
```

The following flags are supported:

| Name            | Description | Type | Default value | Environment variables |
|-----------------|-------------|------|:-------------:|:---------------------:|
| `--help` (`-h`) | show help   | bool |    `false`    |        *none*         |

### `bucket help` subcommand (aliases: `h`)

Shows a list of commands or help for one command.

Usage:

```bash
$ docs [GLOBAL FLAGS] bucket help [command]
```

### `mb` command

Bucket Create.

Usage:

```bash
$ docs [GLOBAL FLAGS] mb [COMMAND FLAGS] [ARGUMENTS...]
```

The following flags are supported:

| Name            | Description | Type | Default value | Environment variables |
|-----------------|-------------|------|:-------------:|:---------------------:|
| `--object-lock` |             | bool |    `false`    |        *none*         |
| `--help` (`-h`) | show help   | bool |    `false`    |        *none*         |

### `mb help` subcommand (aliases: `h`)

Shows a list of commands or help for one command.

Usage:

```bash
$ docs [GLOBAL FLAGS] mb help [command]
```

### `rb` command

Bucket Remove.

Usage:

```bash
$ docs [GLOBAL FLAGS] rb [COMMAND FLAGS] [ARGUMENTS...]
```

The following flags are supported:

| Name            | Description | Type | Default value | Environment variables |
|-----------------|-------------|------|:-------------:|:---------------------:|
| `--help` (`-h`) | show help   | bool |    `false`    |        *none*         |

### `rb help` subcommand (aliases: `h`)

Shows a list of commands or help for one command.

Usage:

```bash
$ docs [GLOBAL FLAGS] rb help [command]
```

### `size` command

Bucket Size.

Usage:

```bash
$ docs [GLOBAL FLAGS] size [COMMAND FLAGS] [ARGUMENTS...]
```

The following flags are supported:

| Name              | Description | Type   | Default value | Environment variables |
|-------------------|-------------|--------|:-------------:|:---------------------:|
| `--delimiter="…"` |             | string |     `"/"`     |        *none*         |
| `--help` (`-h`)   | show help   | bool   |    `false`    |        *none*         |

### `size help` subcommand (aliases: `h`)

Shows a list of commands or help for one command.

Usage:

```bash
$ docs [GLOBAL FLAGS] size help [command]
```

### `policy` command

Bucket Policy.

Usage:

```bash
$ docs [GLOBAL FLAGS] policy [COMMAND FLAGS] [ARGUMENTS...]
```

The following flags are supported:

| Name            | Description | Type | Default value | Environment variables |
|-----------------|-------------|------|:-------------:|:---------------------:|
| `--help` (`-h`) | show help   | bool |    `false`    |        *none*         |

### `policy get` subcommand

Usage:

```bash
$ docs [GLOBAL FLAGS] policy get [COMMAND FLAGS] [ARGUMENTS...]
```

The following flags are supported:

| Name            | Description | Type | Default value | Environment variables |
|-----------------|-------------|------|:-------------:|:---------------------:|
| `--help` (`-h`) | show help   | bool |    `false`    |        *none*         |

### `policy get help` subcommand (aliases: `h`)

Shows a list of commands or help for one command.

Usage:

```bash
$ docs [GLOBAL FLAGS] policy get help [command]
```

### `policy put` subcommand

Usage:

```bash
$ docs [GLOBAL FLAGS] policy put [COMMAND FLAGS] [ARGUMENTS...]
```

The following flags are supported:

| Name            | Description | Type | Default value | Environment variables |
|-----------------|-------------|------|:-------------:|:---------------------:|
| `--help` (`-h`) | show help   | bool |    `false`    |        *none*         |

### `policy put help` subcommand (aliases: `h`)

Shows a list of commands or help for one command.

Usage:

```bash
$ docs [GLOBAL FLAGS] policy put help [command]
```

### `policy help` subcommand (aliases: `h`)

Shows a list of commands or help for one command.

Usage:

```bash
$ docs [GLOBAL FLAGS] policy help [command]
```

### `versioning` command

Bucket Versioning.

Usage:

```bash
$ docs [GLOBAL FLAGS] versioning [COMMAND FLAGS] [ARGUMENTS...]
```

The following flags are supported:

| Name            | Description | Type | Default value | Environment variables |
|-----------------|-------------|------|:-------------:|:---------------------:|
| `--help` (`-h`) | show help   | bool |    `false`    |        *none*         |

### `versioning get` subcommand

Usage:

```bash
$ docs [GLOBAL FLAGS] versioning get [COMMAND FLAGS] [ARGUMENTS...]
```

The following flags are supported:

| Name            | Description | Type | Default value | Environment variables |
|-----------------|-------------|------|:-------------:|:---------------------:|
| `--help` (`-h`) | show help   | bool |    `false`    |        *none*         |

### `versioning get help` subcommand (aliases: `h`)

Shows a list of commands or help for one command.

Usage:

```bash
$ docs [GLOBAL FLAGS] versioning get help [command]
```

### `versioning put` subcommand

Usage:

```bash
$ docs [GLOBAL FLAGS] versioning put [COMMAND FLAGS] [ARGUMENTS...]
```

The following flags are supported:

| Name            | Description | Type | Default value | Environment variables |
|-----------------|-------------|------|:-------------:|:---------------------:|
| `--help` (`-h`) | show help   | bool |    `false`    |        *none*         |

### `versioning put help` subcommand (aliases: `h`)

Shows a list of commands or help for one command.

Usage:

```bash
$ docs [GLOBAL FLAGS] versioning put help [command]
```

### `versioning help` subcommand (aliases: `h`)

Shows a list of commands or help for one command.

Usage:

```bash
$ docs [GLOBAL FLAGS] versioning help [command]
```

### `object-lock` command

Bucket Object Locking.

Usage:

```bash
$ docs [GLOBAL FLAGS] object-lock [COMMAND FLAGS] [ARGUMENTS...]
```

The following flags are supported:

| Name            | Description | Type | Default value | Environment variables |
|-----------------|-------------|------|:-------------:|:---------------------:|
| `--help` (`-h`) | show help   | bool |    `false`    |        *none*         |

### `object-lock get` subcommand

Usage:

```bash
$ docs [GLOBAL FLAGS] object-lock get [COMMAND FLAGS] [ARGUMENTS...]
```

The following flags are supported:

| Name            | Description | Type | Default value | Environment variables |
|-----------------|-------------|------|:-------------:|:---------------------:|
| `--help` (`-h`) | show help   | bool |    `false`    |        *none*         |

### `object-lock get help` subcommand (aliases: `h`)

Shows a list of commands or help for one command.

Usage:

```bash
$ docs [GLOBAL FLAGS] object-lock get help [command]
```

### `object-lock put` subcommand

Usage:

```bash
$ docs [GLOBAL FLAGS] object-lock put [COMMAND FLAGS] [ARGUMENTS...]
```

The following flags are supported:

| Name            | Description | Type | Default value | Environment variables |
|-----------------|-------------|------|:-------------:|:---------------------:|
| `--help` (`-h`) | show help   | bool |    `false`    |        *none*         |

### `object-lock put help` subcommand (aliases: `h`)

Shows a list of commands or help for one command.

Usage:

```bash
$ docs [GLOBAL FLAGS] object-lock put help [command]
```

### `object-lock help` subcommand (aliases: `h`)

Shows a list of commands or help for one command.

Usage:

```bash
$ docs [GLOBAL FLAGS] object-lock help [command]
```

### `lifecycle` command

Bucket Lifecycle.

Usage:

```bash
$ docs [GLOBAL FLAGS] lifecycle [COMMAND FLAGS] [ARGUMENTS...]
```

The following flags are supported:

| Name            | Description | Type | Default value | Environment variables |
|-----------------|-------------|------|:-------------:|:---------------------:|
| `--help` (`-h`) | show help   | bool |    `false`    |        *none*         |

### `lifecycle get` subcommand

Usage:

```bash
$ docs [GLOBAL FLAGS] lifecycle get [COMMAND FLAGS] [ARGUMENTS...]
```

The following flags are supported:

| Name            | Description | Type | Default value | Environment variables |
|-----------------|-------------|------|:-------------:|:---------------------:|
| `--help` (`-h`) | show help   | bool |    `false`    |        *none*         |

### `lifecycle get help` subcommand (aliases: `h`)

Shows a list of commands or help for one command.

Usage:

```bash
$ docs [GLOBAL FLAGS] lifecycle get help [command]
```

### `lifecycle put` subcommand

Usage:

```bash
$ docs [GLOBAL FLAGS] lifecycle put [COMMAND FLAGS] [ARGUMENTS...]
```

The following flags are supported:

| Name            | Description | Type | Default value | Environment variables |
|-----------------|-------------|------|:-------------:|:---------------------:|
| `--help` (`-h`) | show help   | bool |    `false`    |        *none*         |

### `lifecycle put help` subcommand (aliases: `h`)

Shows a list of commands or help for one command.

Usage:

```bash
$ docs [GLOBAL FLAGS] lifecycle put help [command]
```

### `lifecycle rm` subcommand

Usage:

```bash
$ docs [GLOBAL FLAGS] lifecycle rm [COMMAND FLAGS] [ARGUMENTS...]
```

The following flags are supported:

| Name            | Description | Type | Default value | Environment variables |
|-----------------|-------------|------|:-------------:|:---------------------:|
| `--help` (`-h`) | show help   | bool |    `false`    |        *none*         |

### `lifecycle rm help` subcommand (aliases: `h`)

Shows a list of commands or help for one command.

Usage:

```bash
$ docs [GLOBAL FLAGS] lifecycle rm help [command]
```

### `lifecycle help` subcommand (aliases: `h`)

Shows a list of commands or help for one command.

Usage:

```bash
$ docs [GLOBAL FLAGS] lifecycle help [command]
```

### `cors` command

Bucket CORS.

Usage:

```bash
$ docs [GLOBAL FLAGS] cors [COMMAND FLAGS] [ARGUMENTS...]
```

The following flags are supported:

| Name            | Description | Type | Default value | Environment variables |
|-----------------|-------------|------|:-------------:|:---------------------:|
| `--help` (`-h`) | show help   | bool |    `false`    |        *none*         |

### `cors get` subcommand

Usage:

```bash
$ docs [GLOBAL FLAGS] cors get [COMMAND FLAGS] [ARGUMENTS...]
```

The following flags are supported:

| Name            | Description | Type | Default value | Environment variables |
|-----------------|-------------|------|:-------------:|:---------------------:|
| `--help` (`-h`) | show help   | bool |    `false`    |        *none*         |

### `cors get help` subcommand (aliases: `h`)

Shows a list of commands or help for one command.

Usage:

```bash
$ docs [GLOBAL FLAGS] cors get help [command]
```

### `cors put` subcommand

Usage:

```bash
$ docs [GLOBAL FLAGS] cors put [COMMAND FLAGS] [ARGUMENTS...]
```

The following flags are supported:

| Name            | Description | Type | Default value | Environment variables |
|-----------------|-------------|------|:-------------:|:---------------------:|
| `--help` (`-h`) | show help   | bool |    `false`    |        *none*         |

### `cors put help` subcommand (aliases: `h`)

Shows a list of commands or help for one command.

Usage:

```bash
$ docs [GLOBAL FLAGS] cors put help [command]
```

### `cors rm` subcommand

Usage:

```bash
$ docs [GLOBAL FLAGS] cors rm [COMMAND FLAGS] [ARGUMENTS...]
```

The following flags are supported:

| Name            | Description | Type | Default value | Environment variables |
|-----------------|-------------|------|:-------------:|:---------------------:|
| `--help` (`-h`) | show help   | bool |    `false`    |        *none*         |

### `cors rm help` subcommand (aliases: `h`)

Shows a list of commands or help for one command.

Usage:

```bash
$ docs [GLOBAL FLAGS] cors rm help [command]
```

### `cors help` subcommand (aliases: `h`)

Shows a list of commands or help for one command.

Usage:

```bash
$ docs [GLOBAL FLAGS] cors help [command]
```

### `tag` command

Bucket Tagging.

Usage:

```bash
$ docs [GLOBAL FLAGS] tag [COMMAND FLAGS] [ARGUMENTS...]
```

The following flags are supported:

| Name            | Description | Type | Default value | Environment variables |
|-----------------|-------------|------|:-------------:|:---------------------:|
| `--help` (`-h`) | show help   | bool |    `false`    |        *none*         |

### `tag get` subcommand

Usage:

```bash
$ docs [GLOBAL FLAGS] tag get [COMMAND FLAGS] [ARGUMENTS...]
```

The following flags are supported:

| Name            | Description | Type | Default value | Environment variables |
|-----------------|-------------|------|:-------------:|:---------------------:|
| `--help` (`-h`) | show help   | bool |    `false`    |        *none*         |

### `tag get help` subcommand (aliases: `h`)

Shows a list of commands or help for one command.

Usage:

```bash
$ docs [GLOBAL FLAGS] tag get help [command]
```

### `tag help` subcommand (aliases: `h`)

Shows a list of commands or help for one command.

Usage:

```bash
$ docs [GLOBAL FLAGS] tag help [command]
```

### `multiparts` command

Multipart Uploads.

Usage:

```bash
$ docs [GLOBAL FLAGS] multiparts [COMMAND FLAGS] [ARGUMENTS...]
```

The following flags are supported:

| Name            | Description | Type | Default value | Environment variables |
|-----------------|-------------|------|:-------------:|:---------------------:|
| `--help` (`-h`) | show help   | bool |    `false`    |        *none*         |

### `multiparts ls` subcommand

Usage:

```bash
$ docs [GLOBAL FLAGS] multiparts ls [COMMAND FLAGS] [ARGUMENTS...]
```

The following flags are supported:

| Name              | Description | Type   | Default value | Environment variables |
|-------------------|-------------|--------|:-------------:|:---------------------:|
| `--prefix="…"`    |             | string |               |        *none*         |
| `--delimiter="…"` |             | string |     `"/"`     |        *none*         |
| `--help` (`-h`)   | show help   | bool   |    `false`    |        *none*         |

### `multiparts ls help` subcommand (aliases: `h`)

Shows a list of commands or help for one command.

Usage:

```bash
$ docs [GLOBAL FLAGS] multiparts ls help [command]
```

### `multiparts rm` subcommand

Usage:

```bash
$ docs [GLOBAL FLAGS] multiparts rm [COMMAND FLAGS] [ARGUMENTS...]
```

The following flags are supported:

| Name              | Description | Type   | Default value | Environment variables |
|-------------------|-------------|--------|:-------------:|:---------------------:|
| `--key="…"`       |             | string |               |        *none*         |
| `--upload-id="…"` |             | string |               |        *none*         |
| `--help` (`-h`)   | show help   | bool   |    `false`    |        *none*         |

### `multiparts rm help` subcommand (aliases: `h`)

Shows a list of commands or help for one command.

Usage:

```bash
$ docs [GLOBAL FLAGS] multiparts rm help [command]
```

### `multiparts help` subcommand (aliases: `h`)

Shows a list of commands or help for one command.

Usage:

```bash
$ docs [GLOBAL FLAGS] multiparts help [command]
```

### `parts` command

Multipart Parts.

Usage:

```bash
$ docs [GLOBAL FLAGS] parts [COMMAND FLAGS] [ARGUMENTS...]
```

The following flags are supported:

| Name              | Description | Type   | Default value | Environment variables |
|-------------------|-------------|--------|:-------------:|:---------------------:|
| `--key="…"`       |             | string |               |        *none*         |
| `--upload-id="…"` |             | string |               |        *none*         |
| `--help` (`-h`)   | show help   | bool   |    `false`    |        *none*         |

### `parts ls` subcommand

Usage:

```bash
$ docs [GLOBAL FLAGS] parts ls [COMMAND FLAGS] [ARGUMENTS...]
```

The following flags are supported:

| Name            | Description | Type | Default value | Environment variables |
|-----------------|-------------|------|:-------------:|:---------------------:|
| `--help` (`-h`) | show help   | bool |    `false`    |        *none*         |

### `parts ls help` subcommand (aliases: `h`)

Shows a list of commands or help for one command.

Usage:

```bash
$ docs [GLOBAL FLAGS] parts ls help [command]
```

### `parts help` subcommand (aliases: `h`)

Shows a list of commands or help for one command.

Usage:

```bash
$ docs [GLOBAL FLAGS] parts help [command]
```

### `ls` command

Object List.

Usage:

```bash
$ docs [GLOBAL FLAGS] ls [COMMAND FLAGS] [ARGUMENTS...]
```

The following flags are supported:

| Name              | Description | Type   | Default value | Environment variables |
|-------------------|-------------|--------|:-------------:|:---------------------:|
| `--delimiter="…"` |             | string |     `"/"`     |        *none*         |
| `--json`          |             | bool   |    `false`    |        *none*         |
| `--help` (`-h`)   | show help   | bool   |    `false`    |        *none*         |

### `ls help` subcommand (aliases: `h`)

Shows a list of commands or help for one command.

Usage:

```bash
$ docs [GLOBAL FLAGS] ls help [command]
```

### `head` command

Object Head.

Usage:

```bash
$ docs [GLOBAL FLAGS] head [COMMAND FLAGS] [ARGUMENTS...]
```

The following flags are supported:

| Name            | Description | Type | Default value | Environment variables |
|-----------------|-------------|------|:-------------:|:---------------------:|
| `--help` (`-h`) | show help   | bool |    `false`    |        *none*         |

### `head help` subcommand (aliases: `h`)

Shows a list of commands or help for one command.

Usage:

```bash
$ docs [GLOBAL FLAGS] head help [command]
```

### `get` command

Object Download.

Get a single object or add the delimiter (e.g. '/') as path suffix to download recursively.

Usage:

```bash
$ docs [GLOBAL FLAGS] get [COMMAND FLAGS] [ARGUMENTS...]
```

The following flags are supported:

| Name                        | Description                                                            | Type   | Default value | Environment variables |
|-----------------------------|------------------------------------------------------------------------|--------|:-------------:|:---------------------:|
| `--delimiter="…"`           |                                                                        | string |     `"/"`     |        *none*         |
| `--sse-c-key="…"`           | 32 bytes key                                                           | string |               |        *none*         |
| `--sse-c-algorithm="…"`     |                                                                        | string |  `"AES256"`   |        *none*         |
| `--concurrency="…"`         |                                                                        | int    |      `5`      |        *none*         |
| `--part-size="…"`           |                                                                        | int    |      `0`      |        *none*         |
| `--version-id="…"`          |                                                                        | string |               |        *none*         |
| `--range="…"`               | bytes=BeginByte-EndByte, e.g. 'bytes=0-500' to get the first 501 bytes | string |               |        *none*         |
| `--part-number="…"`         |                                                                        | int    |      `0`      |        *none*         |
| `--if-match="…"`            |                                                                        | string |               |        *none*         |
| `--if-none-match="…"`       |                                                                        | string |               |        *none*         |
| `--if-modified-since="…"`   |                                                                        | time   |               |        *none*         |
| `--if-unmodified-since="…"` |                                                                        | time   |               |        *none*         |
| `--help` (`-h`)             | show help                                                              | bool   |    `false`    |        *none*         |

### `get help` subcommand (aliases: `h`)

Shows a list of commands or help for one command.

Usage:

```bash
$ docs [GLOBAL FLAGS] get help [command]
```

### `put` command

Object Upload.

Usage:

```bash
$ docs [GLOBAL FLAGS] put [COMMAND FLAGS] [ARGUMENTS...]
```

The following flags are supported:

| Name                         | Description                                         | Type   | Default value | Environment variables |
|------------------------------|-----------------------------------------------------|--------|:-------------:|:---------------------:|
| `--sse-c-key="…"`            | 32 bytes key                                        | string |               |        *none*         |
| `--sse-c-algorithm="…"`      |                                                     | string |  `"AES256"`   |        *none*         |
| `--part-size="…"`            |                                                     | int    |      `0`      |        *none*         |
| `--concurrency="…"`          |                                                     | int    |      `5`      |        *none*         |
| `--leave-parts-on-error="…"` |                                                     | int    |      `0`      |        *none*         |
| `--max-parts="…"`            |                                                     | int    |      `0`      |        *none*         |
| `--target="…"`               | target key for single file or prefix multiple files | string |               |        *none*         |
| `--acl="…"`                  | e.g. 'public-read'                                  | string |               |        *none*         |
| `--help` (`-h`)              | show help                                           | bool   |    `false`    |        *none*         |

### `put help` subcommand (aliases: `h`)

Shows a list of commands or help for one command.

Usage:

```bash
$ docs [GLOBAL FLAGS] put help [command]
```

### `rm` command

Object Remove.

Remove a single object or add the delimiter (e.g. '/') as path suffix to remove recursively.

Usage:

```bash
$ docs [GLOBAL FLAGS] rm [COMMAND FLAGS] [ARGUMENTS...]
```

The following flags are supported:

| Name                | Description | Type   | Default value | Environment variables |
|---------------------|-------------|--------|:-------------:|:---------------------:|
| `--delimiter="…"`   |             | string |     `"/"`     |        *none*         |
| `--force`           |             | bool   |    `false`    |        *none*         |
| `--concurrency="…"` |             | int    |      `5`      |        *none*         |
| `--dry-run`         |             | bool   |    `false`    |        *none*         |
| `--help` (`-h`)     | show help   | bool   |    `false`    |        *none*         |

### `rm help` subcommand (aliases: `h`)

Shows a list of commands or help for one command.

Usage:

```bash
$ docs [GLOBAL FLAGS] rm help [command]
```

### `cp` command

Object Server Side Copy.

Usage:

```bash
$ docs [GLOBAL FLAGS] cp [COMMAND FLAGS] [ARGUMENTS...]
```

The following flags are supported:

| Name                    | Description                                           | Type   | Default value | Environment variables |
|-------------------------|-------------------------------------------------------|--------|:-------------:|:---------------------:|
| `--src-bucket="…"`      | Source bucket                                         | string |               |        *none*         |
| `--src-key="…"`         | Source key                                            | string |               |        *none*         |
| `--dst-bucket="…"`      | Destinaton bucket                                     | string |               |        *none*         |
| `--dst-key="…"`         | Destination key. When empty, the src-key will be used | string |               |        *none*         |
| `--sse-c-key="…"`       | 32 bytes key                                          | string |               |        *none*         |
| `--sse-c-algorithm="…"` |                                                       | string |  `"AES256"`   |        *none*         |
| `--help` (`-h`)         | show help                                             | bool   |    `false`    |        *none*         |

### `cp help` subcommand (aliases: `h`)

Shows a list of commands or help for one command.

Usage:

```bash
$ docs [GLOBAL FLAGS] cp help [command]
```

### `versions` command

Object Versions.

Usage:

```bash
$ docs [GLOBAL FLAGS] versions [COMMAND FLAGS] [ARGUMENTS...]
```

The following flags are supported:

| Name              | Description | Type   | Default value | Environment variables |
|-------------------|-------------|--------|:-------------:|:---------------------:|
| `--delimiter="…"` |             | string |     `"/"`     |        *none*         |
| `--json`          |             | bool   |    `false`    |        *none*         |
| `--help` (`-h`)   | show help   | bool   |    `false`    |        *none*         |

### `versions help` subcommand (aliases: `h`)

Shows a list of commands or help for one command.

Usage:

```bash
$ docs [GLOBAL FLAGS] versions help [command]
```

### `acl` command

Object ACL.

Usage:

```bash
$ docs [GLOBAL FLAGS] acl [COMMAND FLAGS] [ARGUMENTS...]
```

The following flags are supported:

| Name            | Description | Type | Default value | Environment variables |
|-----------------|-------------|------|:-------------:|:---------------------:|
| `--help` (`-h`) | show help   | bool |    `false`    |        *none*         |

### `acl get` subcommand

Usage:

```bash
$ docs [GLOBAL FLAGS] acl get [COMMAND FLAGS] [ARGUMENTS...]
```

The following flags are supported:

| Name               | Description | Type   | Default value | Environment variables |
|--------------------|-------------|--------|:-------------:|:---------------------:|
| `--version-id="…"` |             | string |               |        *none*         |
| `--help` (`-h`)    | show help   | bool   |    `false`    |        *none*         |

### `acl get help` subcommand (aliases: `h`)

Shows a list of commands or help for one command.

Usage:

```bash
$ docs [GLOBAL FLAGS] acl get help [command]
```

### `acl help` subcommand (aliases: `h`)

Shows a list of commands or help for one command.

Usage:

```bash
$ docs [GLOBAL FLAGS] acl help [command]
```

### `presign` command

Object pre-signed URL.

Usage:

```bash
$ docs [GLOBAL FLAGS] presign [COMMAND FLAGS] [ARGUMENTS...]
```

The following flags are supported:

| Name               | Description | Type     | Default value | Environment variables |
|--------------------|-------------|----------|:-------------:|:---------------------:|
| `--expires-in="…"` |             | duration |     `0s`      |        *none*         |
| `--help` (`-h`)    | show help   | bool     |    `false`    |        *none*         |

### `presign get` subcommand

Presigned URL for a GET request.

Usage:

```bash
$ docs [GLOBAL FLAGS] presign get [COMMAND FLAGS] [ARGUMENTS...]
```

The following flags are supported:

| Name            | Description | Type | Default value | Environment variables |
|-----------------|-------------|------|:-------------:|:---------------------:|
| `--help` (`-h`) | show help   | bool |    `false`    |        *none*         |

### `presign get help` subcommand (aliases: `h`)

Shows a list of commands or help for one command.

Usage:

```bash
$ docs [GLOBAL FLAGS] presign get help [command]
```

### `presign put` subcommand

Presigned URL for a PUT request.

Usage:

```bash
$ docs [GLOBAL FLAGS] presign put [COMMAND FLAGS] [ARGUMENTS...]
```

The following flags are supported:

| Name            | Description | Type | Default value | Environment variables |
|-----------------|-------------|------|:-------------:|:---------------------:|
| `--help` (`-h`) | show help   | bool |    `false`    |        *none*         |

### `presign put help` subcommand (aliases: `h`)

Shows a list of commands or help for one command.

Usage:

```bash
$ docs [GLOBAL FLAGS] presign put help [command]
```

### `presign help` subcommand (aliases: `h`)

Shows a list of commands or help for one command.

Usage:

```bash
$ docs [GLOBAL FLAGS] presign help [command]
```

### `completion` command

Output shell completion script for bash, zsh, fish, or Powershell.

Output shell completion script for bash, zsh, fish, or Powershell. Source the output to enable completion.  # .bashrc source <(sss completion bash)  # .zshrc source <(sss completion zsh)  # fish sss completion fish > ~/.config/fish/completions/sss.fish  # Powershell Output the script to path/to/autocomplete/sss.ps1 an run it.

Usage:

```bash
$ docs [GLOBAL FLAGS] completion [COMMAND FLAGS] [ARGUMENTS...]
```

The following flags are supported:

| Name            | Description | Type | Default value | Environment variables |
|-----------------|-------------|------|:-------------:|:---------------------:|
| `--help` (`-h`) | show help   | bool |    `false`    |        *none*         |

### `completion help` subcommand (aliases: `h`)

Shows a list of commands or help for one command.

Usage:

```bash
$ docs [GLOBAL FLAGS] completion help [command]
```
