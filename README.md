# gitlabclone
Visit groups and subgroups in Gitlab group.  
   
You need to have a valid SSH key in your user directory that is updloaded in Gitlab.  
You need an API token with read_repository and read_api right.  
  
# Usage
## #NAME:
   gitlabclone - clone all the projects and subprojects below a group or project

## USAGE:   
   gitlabclone [options]... -r domain  -d destination
  
## COMMANDS:
   help, h  Shows a list of commands or help for one command  
  
## OPTIONS:
   -t token  
   --token  _gitlab access token_  the access token for accessing the git gitlab repository through https  

   -r domain    
   --domain _gitlab domain_ root domain of the gitlab repository (default endpoint will then be https://domain/api/v4) 
 
   -a api/v4 root path  
   --api-path _gitlab api v4 root path_  root endpoint path for the gitlab domain  
     
   -g groupid  
   --group _group id_ id of gitlab group (if not provided all groups / projects of domain are retrieved)  
     
   -k sshkey relative path  
   --ssh-relative-path _relative path of ssh key_ relative user path for ssh key (i.e. .ssh/id_rsa)  
     
   -d destination    
   --destination _absolute path_ absolute path where to clone the projects  

   -l level      
   --log-level _level_ Log level (can be error/warning/info/debug/trace) (default: "info")  
     
   -h  
   --help show help (default: false)  


