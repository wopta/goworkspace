# Tag Functions script
## Usage
The script expects the following argurments:
- `functions`: a comma separated list of functions to release
- `from`: (optional) source environment. If not given the target environment will be used
- `target`: the target environment

## Example
For new function releases on dev environment
```bash
go run main.go --functions=broker,callback,renew,payment --target=dev
```

For promoting functions from dev to uat
```bash
go run main.go --functions=broker,callback,renew,payment --from=dev --target=uat
```

# Tag Functions Cloudbuild
## Usage
The trigger must be called manually with the needed arguments for the previously mentioned script

The access token for comunicating with git is present in secrets manager and it is generated in the project access token page of the gitlab project