# Cat API - Beego Project

This is a Beego-based API project for managing information about cats. Follow the instructions below to clone, set up, and run this project on your local machine in both Windows and Linux environments.

## Prerequisites

Before starting, ensure you have the following installed:

- [Go](https://golang.org/doc/install) (version 1.18 or later)
- [Git](https://git-scm.com/)

---

## Installation

### Step 1: Install Beego Framework
```bash
go install github.com/beego/beego/v2@latest
bee version
```

### Step 2: Check the GOPATH
```bash
go env GOPATH
```
### Step 3: Set GOBIN to the path variable
```bash
setx PATH "%PATH%;%GOPATH%\bin" # on windows
# on linux
nano ~/.bashrc
#or
nano ~/.zshrc
#Add the following line at the end of the file:
export PATH=$PATH:$GOPATH/bin
#Save the file and reload the configuration:
nano ~/.bashrc
#or
nano ~/.zshrc
#Verify the PATH by running:
echo $PATH
```
### Step 4: Make the Src folder:
```bash
#Navigate to GOPATH
cd %GOPATH% #on windows
cd $GOPATH #on linux
#create the directory
mkdir src
cd src
```
### Step 5: Clone the Repository

Navigate to src directory and run

```bash
git clone https://github.com/Marjia029/cat-api.git
cd cat-api
```

### Step 6: Run The Application
```bash
go mod tidy
bee run
```
Now go to http://localhost:8080 and check the application 

## Testing
Open the Terminal and Run
```bash
 go test ./tests -v
 ```
 For test Coverage:
 ```bash
 go test -coverprofile=coverage.out ./...
 go tool cover -html coverage.out
 ```
 Now check the dropdown to see the coverage percentage of every file.

