gocp.go V2
	
	Be used to copy any FILES or DIRS or both to a DIR or FILE.
	Upgrade from gocp V1.
	Other than V1, you don't need double quotes to surround the glob pattern string,
	as the shell will expansion the glob pattern string to multi sources files or directorys,
	and we make the last argument of shell as the destination.
	eg: gocp test* /tmp
	    Every file or directory that matched by "test*" will be a source, and "/tmp" will be
	    destination.

	YOU CAN DO:
		1. FILE to NEW_FILE;
		2. [DIRs, FILEs, ...] to NEW_DIR;
		3. FILE to EXIST_FILE (overwrite);
		4. [DIRs, FILEs, ...] to EXIST_DIR;

