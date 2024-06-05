# **Tutorial**:

------------------------------------------------------------------------------------------------------------
1) Make sure to download Go on version 1.22.2 based on OS: https://go.dev/dl/
2) Download the dependencies using: go mod download

------------------------------------------------------------------------------------------------------------
3Run one of these commands depending on operating system and architecture,
   if the executables are not in the bin folder

**LINUX** <br>
GOOS=linux GOARCH=amd64 go build -o bin/app-amd64.exe main.go

**Mac (m1):** <br>
GOOS=darwin GOARCH=arm64 go build -o bin/app-arm64-darwin main.go

------------------------------------------------------------------------------------------------------------

4After one of the executable files are in the bin, add the input files in the same folder level as main.go.


**Linux:** <br>
.\bin/app-amd64.exe < -file1 init-dp.txt -file2 input-dp.txt > output.txt

**Mac (m1):** <br>
./bin/app-arm64-darwin -file1 init-dp.txt -file2 input-dp.txt > output.txt




_****if the name is different than the init-dp.txt and input-dp.txt adjust the name accordingly, <br>
*but make sure that the first file is always the init file******_

Run the relative command based on OS

------------------------------------------------------------------------------------------------------------
# **File descriptions**: 

-The main.go contains all my functions and the presentation shell.

-Go.mod keeps track of my projects dependencies similar to package.json

-Go.sum verifies the integrity of downloaded dependencies 

-Myapp.log keeps track of all the logs, used for debugging

-The bin holds the executables