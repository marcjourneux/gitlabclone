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
   --access-token value, -t value  gitlab access token
 
   --group value, -g value         id of gitlab group
 
   --ssh value, -k value           relative user path for ssh key
   
   --destination value, -d value   local path where to clone the projects (folders will be created)
   
   --log-level value, -l value     Log level (error/warning/info/debug/trace) (default: "Info")
   
   --help, -h                      show help (default: false)


