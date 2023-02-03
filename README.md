# gitlabclone
Visit groups and subgroups in Gitlab group 
You need to have a valid SSH key in your user directory that is updloaded in Gitlab.
You need an API token with read_repository and read_api right.

# Usage
NAME:
   gitlabclone - clone all the projects and subprojects below a group or project

USAGE:
   gitlabclone [global options] command [command options] [arguments...]

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --token , -t                     gitlab access token
   --domain , -d                    root domain (default endpoint will be https://domain/api/v4)
   --api-path, -a                   endpoint path for the gitlab domain
   --group , -g                     id of gitlab group (can be empty for retrieving all groups)
   --ssh-relative-path , -k         relative user path for ssh key (i.e. .ssh/id_rsa)
   --destination , -d               local path where to clone the projects (folders will be created)
   --log-level , -l                 Log level (error/warning/info/debug/trace) (default: "Info")
   --help, -h                       show help (default: false)


