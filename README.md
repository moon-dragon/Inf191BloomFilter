<<<<<<< HEAD
# github.com/vlam321/inf191BloomFilter
=======
# github.com/vlam321/Inf191BloomFilter
>>>>>>> upstream/develop

## About
This project is a proof of concept. Our problem assumption is that the current system queries userid:email pairs from a giant database. Though this provides an accurate membership result, the network overhead of querying the database is too large. Therefore, we will build a bloom filter to improve data management performance. The current scope will be to integrate and test an existing bloom filter and expand it to a distributed node implementation. We will also build utilities for testing our bloom filter locally. These tools include a data generator and visualization metrics.

## Contribution
To start contributing to this project, simply fork this repository and clone it using after replacing "vlam321" with your user name, Like so: 
<<<<<<< HEAD
`git clone https://github.com/yourGithubUserName/github.com/vlam321/inf191BloomFilter`
=======
`git clone https://github.com/yourGithubUserName/github.com/vlam321/Inf191BloomFilter`
>>>>>>> upstream/develop

## Dependencies
The following dependencies will be needed to run some of the files in this repo. Use the commands below to install them if needed.
### Go MySQL Driver
` $ go get -u github.com/go-sql-driver/mysql`
### Bloom Filter
` $ go get -u github.com/willf/bloom`
### Testify
` $ go get github.com/stretchr/testify`
### Graphite API client
` $ go get github.com/marpaia/graphite-golang`

## Grafana and Graphite Setup
1. Download and install [Docker Toolbox](https://www.docker.com/products/docker-toolbox)
2. Open Docker Quickstart Terminal
3. Pull the Grafana and Graphite images with the following commands:
```
$ docker pull grafana/grafana
$ docker pull hopsoft/graphite-statsd
```
4. Run the images in the background using the following commands:
```
$ docker run -d --name=grafana -p 3000:3000 grafana/grafana
$ docker run -d --name graphite --restart=always -p 80:80 -p 2003-2004:2003-2004 -p 2023-2024:2023-2024 -p 8125:8125/udp -p 8126:8126 hopsoft/graphite-statsd
```
5. Make sure that Grafana and Graphite are running by going to localhost at port 80 and 3000 (separately!)
6. Run graph.go under cmd/graph to see that it is affecting that the graphs in your Graphite page 

Note: localhost might not work. I'm currently trying to find a way to make it work, but for now, if it doesn't just type `docker-machine ip` in the quickstart terminal to get the ip address. Use that instead of 'localhost' along with the port number to check if the containers are running correctly. Also, if you're running the getting an error while running graph.go, you might have to also change 'localhost' to the ip address you get from docker-machine. There are also some useful commands that you might want to use in `docker --help`, or you can also read more about how to use Docker in the official [doc](https://docs.docker.com/).

## Test Data
### bloomDataGenerator
Test Data can be generated using the bloomDataGenerator package. The function GenData will require 3 arguments: amount of users, and a minimum and maximum values to set the range for the number of email addresses to be generated.
```
data := bloomDataGenerator.GenData(5, 100, 1000)
```
### databaseAccessObj
To populate database with these data, ensure that your database, call the InsertDataSet function and a map[int][]string as instance as argument.
```
dao := databaseAccessObj.New(dsn)
dao.Insert(data)
```
There are also other useful functions in databaseAccessObj.go.
### updateDatabase
Additionally, you can get your data store set up in command line using updateDatabase.go. There are a few ways to do it. You can run it directly, or use go install to create a executable. Either way should work fine as long as you provide the correct command. Here's a few examples on how to use it.
#### Populating the database with tables
```
$ updateDatabase -cmd=mktbls
```
Note: You can run the script without using 'go run' only after you did 'go install' to create the binary file. The command above would create 15 tables in the 'unsubscribed' database, so make sure you have that setup first. You should really start with an empty database (without any tables).
#### Repopulating the database with values
```
$ go run updateDatabase.go -cmd=repopulate -numUser=10 -minEmail=1000 -maxEmail=10000
```
Note: You can use 'go run' to run the script, which will clear the current values in existing unsub tables and add new values in.
#### Adding new values (without deleting existing data)
```
$ go run updateDatabase.go -cmd=add -numUser=2 -minEmail=100 -maxEmail=200
```
