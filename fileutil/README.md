## fileutil
Provides utilities for interacting with files on the filesystem, such as determining if a file is Binary or a Symlink.

# IsBinaryFile
Attempts to determine if the provided file is Binary by reading pieces of it from the filesystem.

# IsSymlink
Determines if the file is a Symlink

# IsDir
Determines if the file is a directory or regular file.

# DirToArray
Collapses a directory tree into a slice of files, including the full path to the file. File filters and Directory filters can be supplied to limit the results, such as file filters that remove Binary or Symlink files from the results. The default file filter exclused symlinks and binary files, and the default directory filter excludes a .git directory. These can be supplied to fit your needs.

```
files := fileutil.DirToArray(dir, false, fileutil.SymlinkFileFilter, fileutil.DefaultDirectoryFilter)
```

# FilterExtBlacklist
FilterExtBlacklist takes a list of extensions and a slice of files and removes any files in that list that have an extension matching an item in the blacklist.

# FilterExtWhitelist
FilterExtWhitelist takes a list of extensions and a slice of files and removes any files in that list that do not have an extension matching an item in the whitelist.