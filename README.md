#OBO

Revolutionizing buying and selling on college campuses.
OBO is a location-based iOS application that lets users post
goods and items for sale, as well as request postings in a 
nearby radius. 

##Features
- Create listings for goods and services
- Create requests
- Edit or remove posts
- Get active posts within your radius on refresh
- Flag posts for spam
- Once a match is made, the buyer and seller are connected


#Installation Instructions

##Installing Go
First, you need to follow these instructions:
https://golang.org/doc/install

The easiest way to install Go is with the following installer:
https://storage.googleapis.com/golang/go1.4.2.darwin-amd64-osx10.8.pkg

If it moves, or we change versions, or something... I'll update as approprate. For consistency, let's all go with v1.4.2, for now.

##Understand the Go Workspace
Before you do anything, read/skim this:
https://golang.org/doc/code.html

The gist:
You will have ONE Go workspace for EVERY Go project you ever write. Given the group's resistance to Go in the first place, I wouldn't be surprised if this is the only Go project you ever write. No matter. If you've installed correctly, the following environment variable will be set:

GOPATH -> this is the root directory of your go workspace

Check to make sure that your GOPATH exists by typing the following command:
```
// prints the GOPATH
printenv GOPATH
```
Next, check to make sure that you can run `go` from the command line. Simply type `go`, and if it finds the command, you're set.

If either the varible isn't set appropriately, or the go command isn't found, modify our bash profile with the following (generally found at `~ .bash_profile):

```
// Adds the go command to your path
export PATH=$PATH:/usr/local/go/bin
// Creates the GOPATH. If you don't want the go workspace in your home
// directory, you can change the GOPATH here
export GOPATH=$HOME/go
```

Finally, add the following to your bash profile:
```
// This will allow you to write programs, use the 'go install' command,
// and run the program directly by typing the name of the program
export PATH=$PATH:$GOPATH/bin
```

### The Go Workspace:
This bit is important. Go should generate the file structure for your workspace automatically, but in case it doesn't, the toplevel of the workspace ($GOPATH) should have 3 folders:
- src 	// Source code for projects go here
- pkg	// If you build separate packages (which we prob. wont), they 
		// go here. Otherwise, non-golang packages will be put here
- bin 	// When you build your code with go install, the resulting binary
		// will be placed here, with the name of the 'main' package file
		// as the command name. So if I have a file called `hello.go`,
		// the executable `hello` will be placed here. You can then run
		// `hello` from the command line

Don't worry about the pkg and the bin directories for now. The src directory should be organized in the following fashion:

```
$GOHOME/src/{remote repository domain}/{username}
```

For example, my file tree looks like this:

```
$GOHOME/github.com/beisner/OBO
```

I suggest you structure it the same way, since that's how we'll do it on the server.

Within the source directory for the project (/OBO in this case), code structure is very flat. For now, we won't have any subfolders. That means that all source code and all testing code will be in this file. Go is built to be concise, and the built in testing framework is excellent, so we don't really need much more structure.

#OBO in Go
General information about how to build, run, test, and deploy the project. Also some information about the libraries that we're using.

##Getting the source code
Execute the following command when you're in the '$GOHOME/src/{remote repository domain}/{username}' directory:
```
git clone https://github.com/wcdavis/OBO.git
```

##How to Code in Go
You need to learn go on your own.

##Building OBO
There's a makefile in there, so all you really need to do is run
```
make backend
```
and you'll be set.

##Running OBO
We're going to want to pass a few arguments and files to the go program, so there's just run the following command and the rest will be taken care of in the makefile:

```
make run-backend
```

##Testing OBO
Go has a really nice testing framework. Basically, all you do is write a testfile that has testing functions, include the 'testing' package, and name the file *_test.go, and if you run `go test` in the code directory all the tests will be performed, with a nice little output that tells you if you've passed the test. The testing is automated in the makefile, with the following command:
```
make test-backend
```

##Deploying OBO
Deployment is automated. The deployment script will run all tests, and will only deploy if all the tests pass. There will be 3 servers running: test, stage, and live. This may sound a bit excessive, but it won't cost us much (if anything) extra in terms of money, and will save us a lot of confusion when developing applications against the backend. Here's what they're for:

- **test:** So you want to try out your swanky new feature on a remote server? No problem! Ping the other developers to make sure that nobody is testing the backend currently, and deploy away!

- **stage:** Have we, as a team, decided that we've reached a stable point in development, and want to test a release candidate? Awesome! Good work! Deploy to stage, so the iOS developers can test thoroughly against the backend. iOS developers should always be coding against stage. There will be a bunch of test data and users on live.

- **live:** Once stage has been thoroughly tested, and we're ready to put our release candidate into the wild, ***BE ABSOLUTELY SURE*** you're ready to deal with whatever might happen, and deploy to live. THIS IS NOT FOR TESTING. live is for the current release. This might seem not-so-important right now, but if we release once and get some users, then we release a broken version to live, our users will leave us we have the plague. Really, don't do this unless you are dead certain.

The following make commands will do the trick:

```
// deploy to test
make deploy

// deploy to stage
make deploy-stage

// deploy to live
make deploy-live-definitely-stable

// I'm not joking about that last one. The length is a deterrent.
```

##Packages, frameworks, etc.
We'll be using a few different frameworks for this project. Here is the current list:

- [github.com/gorilla/mux](http://www.gorillatoolkit.org/pkg/mux) : routing library for easily adding paths and accessing url path variables
- [emicklei/go-restful](https://github.com/emicklei/go-restful) (maybe) : REST library, pretty heavyweight but seems to do a reasonably good job at being simple. I only like this one over the golang net/http package because it has Skritter integration, with would allow us to automatically generate documentation for the API (whenever you add a new endpoint, it'll automatically fill in the approprate inforomation about the API)

To resolve all the dependencies, all you need to do is to execute the following command in the source directory:

```
go get
```

##Project Structure
Design forthcoming.

#Git Etiquitte
We need this.

##Using Git
Learn how to use git on your own. Seriously, super useful stuff.

##Branching
If you want to make any sort of change to the code, ***CREATE A NEW BRANCH, AND DO NOT MAKE ANY CHANGES ON MASTER***. Just so we're clear on who's doing what and when, I'd like to use the following naming conventions for branches:
```
{type}/{handle}/{branch name}
```

The tree types of branches are as follows:
- f : feature branch, created when making a new feature
- b : bug branch, created when fixing/patching a bug
- i : information, only used when changing the readme or adding non-code content (images, etc) to the project

For example, creating this readme, I branched master into the following branch:

i/beisner/readme_v1

##Committing
Every commit should build. It may fail tests, but please make it build. Also, this goes without saying, but descriptive messages, please.

##Pull Requests
When you think you're ready to merge with master, go to Github and create a pull request with master. Describe what changes you've made. ***DO NOT MERGE THE PULL REQUEST UNTIL YOU ASK A TEAMMATE TO LOOK AT YOUR CODE.*** This is very important. We all have the power to merge, but we should talk about merges and code quality before we make them with master.

##Keeping Code Up To Date
Please, try to pull from master at least every other day. When in your working branch, you can execute the following command to update the code in your branch:
```
git pull origin master
```

##Stable Branch
We will also have a stable branch, which nobody should touch! This is the branch that is currently in production on live, so we know exactly what's running on our servers.