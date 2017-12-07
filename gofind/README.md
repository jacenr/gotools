# gofind V2
	
	## Search file or directory by name or size or modify time in specify paths or current work directory.

	## For example:
		### ./gofindv2 -n=reg,.*test.*
			Search files and direcorys which matched by regular expression ".*test.*" 
			in current work directory.
			Note: "reg" and ".*test.*" should be separated by comma.
		### ./gofindv2 -n=full,nginx -p=/usr
			Search nginx in /usr. Note the comma separator.
		### ./gofindv2 -m=">,20171205000000" -n=sub,test -d=o -p=/tmp,/mnt
			Search file whose name contain "test" and modify time after "Dec 5 00:00:00 2017", 
			and only directory are outputted.
			Note the comma separator and option with digits should be surrounded by double quotes.

	## Options:
		### -n 
			Search file by file name. For example: "full,testName" or "sub,test", or "reg,.*test.*".
			"full": full string matching, "sub": sub string matching, "reg": regular expression matching.
	
		### -s
			Search file by file size. Option example: ">=,1024", unit B.
	
		### -m
			Search file by modify time. Option format: ">=,20171206114930".
	
		### -d
			Filter file or directory by type. "dir": "directory". For example: -d="o", 
			only dierctory name is outputted.
	
		### -p
			Specify the search path, default current work directory. Multi paths separated by comma.


