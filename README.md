# TUTORIAL OF USEING THIS DEMO SERVER

##
1.Assume you have an Unix-like system . And You have intall the Postgresql  and the Go language envirment

2.workplace direcory instructions:

```
bin/
   start.sh             # it is used to start the server  usage : enter the project directory excute: bin/start.sh
                         
etc/
    config.json         # it contains the db configuration and so on . before start the server you should change the config match your environment
src/
    handler/            # the url handler
        relation.go     # deal with the relation request 
        user.go         # deal with the user info request
    model/              # db table model
        model_types.go  # define some data type
        db_init.go      # db init
        relationships.go# table relationships model
        user.go         # table users model
    util/
        config.go       # config init from the json type file 
        log.go          # log util
    main.go             #
log/                    # when the server start user request will genarate the log file : info.log err.log warn.log
                        # the files represent diffent level log
vendor/                 # third party package directory

```


3.use postgresql client to connect pg server then execute the bin/install.sql content

4.configurate your pg server ip, port, username and password.

4.Start the server :
* clone the project to your [GOPATH](https://github.com/golang/go/wiki/GOPATH "Title") (asumme you have install golang envirment)
* execute : bin/start.sh(your local machine 8080 port is avaliable to use)
* use tool to request the server    


TODO 
* circuit breaker


