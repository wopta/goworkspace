# Tag Modules script
## Usage
The script expects the following argurments:
- `modules`: a comma separated list of modules to release
- `modulePath`: (optional) the project name prefix for the modules to be releases. Defaults to github.com/wopta/goworkspace/

## Example
```bash
go run main.go --modules=document
```

# Tag Modules Cloudbuild
## Usage
The trigger must be called manually with the needed arguments for the previously mentioned script

The access token for comunicating with git is present in secrets manager and it is generated in the project access token page of the gitlab project