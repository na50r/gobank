# Compilers
## Windows 
* Issues can occur when working with the sqlite3 package
* What worked was using the TDM-GCC Compiler: https://jmeubank.github.io/tdm-gcc/
* GCC: gcc.exe (MinGW-W64 x86_64-ucrt-posix-seh, built by Brecht Sanders, r8) 13.2.0
* GO: go1.22.5 windows/amd64

## Linux (Ubuntu WSL)
* GCC: gcc (Ubuntu 11.4.0-1ubuntu1~22.04) 11.4.0
* GO: go1.22.5 linux/amd64

## SQLite
* Some operations differ from Postgres, such as how values are injected into the SQL string
* Or how to handle rows operations: https://stackoverflow.com/questions/32479071/sqlite3-error-database-is-locked-in-golang