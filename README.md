# sss

Yet another S3 client.

## Documentation

[DOCS.md](DOCS.md) contains the generate documentation.

## Usage

For shell completion follow the instructions from `sss complete`.

```
NAME:
   sss - S3 client

USAGE:
   sss [global options] [command [command options]]

COMMANDS:
   buckets      
   bucket       
   mb           
   rb           
   multiparts   
   parts        
   ls           
   cp           
   put          
   rm           
   get          
   head         
   presign      
   policy       
   cors         
   object-lock  
   lifecycle    
   versioning   
   acl          
   versions     
   help, h      Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --endpoint string    
   --insecure           
   --region string      
   --path-style         
   --profile string     (default: "default")
   --bucket string      
   --secret-key string  
   --access-key string  
   --verbosity uint     (default: 1)
   --help, -h           show help
```

## Configuraton

`sss` uses the AWS SDK and the same configuration files and environment variables (e.g. `AWS_ACCESS_KEY_ID` and `AWS_SECRET_ACCESS_KEY`) as the AWS CLI.

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
