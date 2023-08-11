# whoisjsonapi domainScan
Script that reads list of domains and retrieve custom info in JSON line format


I find out random freelance job and make a solution for personal practice so I'm 
willing to share it. 

This script scans custom set of data that is defined by needs of client (check structs
that are initialized in start of code) from [whoisjsonapi](https://whoisjsonapi.com/).
If you are in need of additional fields that API can provide you, just adapt structure
of the code (main.go) to a fields that are avaiable on API response. 

Setup of environemt variable is needed for script to be executed so lets explain env 
variables in details (.dScan).

Example of fully configured enviroment variables
```
DSCAN_DOMAIN_FILE = domains.txt
DSCAN_LOG_FILE_NAME = domainScan.log
DSCAN_JSON_FILE_NAME = results
DSCAN_JSON_EXTENSION = log
DSCAN_JSON_SIZE_LIMIT = 500000000
DSCAN_WHOIS_API_TOKEN = Bearer {your API token from whoisjsonapi.com}
```

DOMAIN_FILE - is file that containts all the list of domains that are in formated as below
(each line represent one input and its read line by line):
```
xiaohongshu.com
3dmgame.com
tradingview.com
3dmgame.com
xiaohongshu.com
instructure.com
mediafire.com
alibaba.com
cnki.net
youtube.com
```

LOG_FILE_NAME - this variable is set by your own choice (name + extension) and its
used to track any kind of INFO or ERROR logs that will help you debug or check 
process of code nad its formatted as standard log file with custom message.
Includes timestamp and message of info or error.

JSON_FILE_NAME - represent name of script output

JSON_EXTENSION - extension of script output 

So basically naming will go {JSON_FILE_NAME}-N.{JSON_EXTENSION} where N represents number
that increase after file limit is reached (e.g. results-1.log based on full config of env
variables in .dScan)

JSON_SIZE_LIMIT (Bytes) - is used to specify size of each output file, and it means if
output file reaches its limit, new one is created (500MB in this case from .dScan)

Feel free to contact me if you need some help of code refactoring or suggestion.

Contact: domagojpr@gmail.com

