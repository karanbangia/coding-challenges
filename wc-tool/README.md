
# wc-tool
WC Tool written as per the specifications in John Crickett's Coding Challenges written in Golang. It supports 5 basic functionalities:
- Count Characters in a file
- Count Words in a file
- Count Bytes in a file
- Count Lines in a file
- Count Lines from Standard Input


## How to run the Project
- Fork the repository and download on your local
- Make sure you have the latest version of go downloaded and installed in your system. If not the same can be done from here https://go.dev/doc/install
- Move to the home directory (containing the `main.go` file)
- Build the file `go build -o ccwc`
- Run the program with the command `./ccwc <command line options> <fileName>`
- The Allowed command line options are as follows
    - `-c` : Gives byte count
    - `-l` : Gives line count
    - `-w` : Gives Word count
    - `-m` : Gives char count

## Final Step
- In this step we are able to read from standard input if no filename is specified
``>cat test.txt | ccwc -l``
 ```` 7137 ````