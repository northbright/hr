# update-employees

update-employees is a program which read and parse employee data in given CSV file and update the data stored in Redis database. It's written in [Golang](http://golang.org).

#### Usage

    -a string
    	Redis server address. Ex: -a='127.0.0.1:6379'
    -f string
    	path of CSV file which contains employee data. Ex: -f='employees.csv'
    -p string
    	Redis password. Ex: -p='my_password'

    usage:
        update-employees -a=<Redis server address> -p=<Redis password> -f=<csv file>

    e.g.
        update-employees -a='127.0.0.1:6379' -p='' -f='employees.csv'




