# rabbitmq

bunny tests

## The Environment variables

For simplicity [direnv](https://direnv.net/) is used, and the environment variables used among all the code are in the file `.envrc`

    export RABBITMQ_USER=bunny
    export RABBITMQ_PASS=test
    export RABBITMQ_HOST=my-rabbit
    export RABBITMQ_VHOST=hole

In case you don't have it you just need to declare those variables with your own values

## setting up the go environment

To run them set your go environment https://golang.org/doc/install

For example using `$HOME/go` for your workspace

    $ export GOPATH=$HOME/go

Create the directory:

    $ mkdir -p $HOME/go/src/github.com/nbari

Clone the project into that directory:

    $ git clone https://github.com/nbari/rabbitmq.git $HOME/go/src/github.com/nbari/rabbitmq

## running the tests

Just `cd` into a directory where is a `main.go` and run:

    $ go get  (in case need to get the dependecies)
    $ go run main.go


For example to test creating random queues:

    $ cd $HOME/go/src/github.com/nbari/rabbitmq/ranbom-queues/no-durable
    $ go run main.go

In case need to compile just run, it will create a binary the one could be used from any location, just keep in mind it will need the `ENV` vars

    $ go build
