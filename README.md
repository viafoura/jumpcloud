# JumpCloud

NAME:  
   ClearCare Jumpcloud - Work w/ the Clouds of Jump

USAGE:  
   ClearCare Jumpcloud [global options] command [command options] [arguments...]

VERSION:  
   0.0.0

AUTHOR(S):

COMMANDS:  
   tag          Tag operations  
   system       System operations  
   help, h      Shows a list of commands or help for one command  

GLOBAL OPTIONS:  
   --config, -c "/opt/jc/jcagent.conf"  Specify an alternate agentConfig Default: /opt/jc/jcagent.conf  
   --verbose, -V                        Be verbose  
   --help, -h                           show help  
   --version, -v                        print the version  



Its assumed that this will be run from the system, if not you will need the config file of the system that you are wanting to edit  

Add a tag:  
jumpcloud system addTag  

Remove a tag:  
jumpcloud system removeTag

