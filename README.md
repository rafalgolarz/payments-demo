This project has been tested on Linux Mint 18.3 XFCE 64b.

Altough the code is built to be not tide to a storage type, it has a concrete implementation for MySQL to check the functionality in more real scenario.

If you don’t have a local MySQL database setup, you can use docker (port 3306 should be open and not utililised).


```
go get  github.com/rafalgolarz/payments-demo/cmd/paymentsd

cd $GOPATH:/github.com/rafalgolarz/payments-demo

docker pull mysql/mysql-server

docker run --name=mysql_payments -p3306:3306 -e MYSQL_ROOT_PASSWORD=password -e MYSQL_ALLOW_EMPTY_PASSWORD=yes -d mysql/mysql-server --default-authentication-plugin=mysql_native_password

```

it may take a couple of seconds, run this to check if the mysql server is up:


```
docker logs mysql_payments

```

Let’s create our database first:


```
docker exec -i mysql_payments mysql -uroot -ppassword -e "create database payments_demo"

```

then let’s import tables:

```
docker exec -i mysql_payments mysql -uroot -ppassword payments_demo < ./scripts/payments.sql

```

Let’s do the same for our test db:


```
docker exec -i mysql_payments mysql -uroot -ppassword -e "create database payments_demo_test"

```

then let’s import tables:


```
docker exec -i mysql_payments mysql -uroot -ppassword payments_demo_test < ./scripts/payments_testdb.sql

```

Now let’s jump to the console to check if databases exist:


```
docker exec -ti mysql_payments mysql -uroot -ppassword

```

Obviously it’s not recommended to pass password in the command line (may stay in the history) but for the sake of simplicity, we can do it that way.

Now time to run our API:

```
cd $GOPATH:/github.com/rafalgolarz/payments-demo/cmd/paymentsd/

go build .

./paymentsd
```