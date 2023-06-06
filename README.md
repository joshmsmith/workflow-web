# webapp with Temporal workflow backend
Sample golang web app with local db task queue and Temporal Cloud workflow backend interaction

## start mysql database
The sample mysql database has been configured to run using docker-compose locally and initialise the database with users and sample data.
```
cd mysql
docker-compose up -d
docker logs mysql
```

Sample data:
```
docker exec -it mysql mysql -u root -p
Welcome to the MySQL monitor.  Commands end with ; or \g.
...
mysql> describe dataentry.accounts;
+-----------------+--------------+------+-----+-------------------+-------------------+
| Field           | Type         | Null | Key | Default           | Extra             |
+-----------------+--------------+------+-----+-------------------+-------------------+
| account_id      | int unsigned | NO   | PRI | NULL              | auto_increment    |
| account_number  | int          | NO   |     | NULL              |                   |
| account_name    | varchar(30)  | NO   | UNI | NULL              |                   |
| account_balance | float        | NO   |     | NULL              |                   |
| datestamp       | timestamp    | NO   |     | CURRENT_TIMESTAMP | DEFAULT_GENERATED |
+-----------------+--------------+------+-----+-------------------+-------------------+
5 rows in set (0.01 sec)

mysql> describe moneytransfer.transfer;
+-------------+--------------+------+-----+-------------------+-------------------+
| Field       | Type         | Null | Key | Default           | Extra             |
+-------------+--------------+------+-----+-------------------+-------------------+
| id          | int unsigned | NO   | PRI | NULL              | auto_increment    |
| origin      | varchar(30)  | NO   | MUL | NULL              |                   |
| destination | varchar(30)  | NO   |     | NULL              |                   |
| amount      | float        | NO   |     | NULL              |                   |
| reference   | varchar(40)  | NO   |     | NULL              |                   |
| status      | varchar(30)  | NO   |     | NULL              |                   |
| t_wkfl_id   | varchar(50)  | YES  |     | NULL              |                   |
| t_run_id    | varchar(50)  | YES  |     | NULL              |                   |
| t_taskqueue | varchar(50)  | YES  |     | NULL              |                   |
| t_info      | varchar(250) | YES  |     | NULL              |                   |
| datestamp   | timestamp    | NO   |     | CURRENT_TIMESTAMP | DEFAULT_GENERATED |
+-------------+--------------+------+-----+-------------------+-------------------+
11 rows in set (0.00 sec)

mysql> select id,origin,destination,amount,reference,status from moneytransfer.transfer;
+----+--------+-------------+--------+------------------+-----------+
| id | origin | destination | amount | reference        | status    |
+----+--------+-------------+--------+------------------+-----------+
|  1 | bill   | jim         |    120 | IOU              | REQUESTED |
|  2 | jane   | sally       |    107 | FOOD MONEY       | REQUESTED |
|  3 | ted    | harry       |    100 | CART123          | REQUESTED |
|  4 | bill   | ted         |     10 | transfer request | REQUESTED |
+----+--------+-------------+--------+------------------+-----------+
4 rows in set (0.00 sec)

mysql> quit
Bye
```

## Temporal Cloud configuration
This example assumes that you have a temporal cloud configured and have local client certificate files for your namespace.
The values are passed into the demo app using environment variables, example direnv .envrc file is included in the repo:

```
# Temporal Cloud connection
# region: us-east-1
export TEMPORAL_HOST_URL="myns.abcdf.tmprl.cloud:7233"
export TEMPORAL_NAMESPACE="myns.abcdf"

# tclient-myns client cert
export TEMPORAL_TLS_CERT="/Users/myuser/.temporal/tclient-myns.pem"
export TEMPORAL_TLS_KEY="/Users/myuser/.temporal/tclient-myns.key"

# Optional: path to root server CA cert
export TEMPORAL_SERVER_ROOT_CA_CERT=
# Optional: Server name to use for verifying the server's certificate
export TEMPORAL_SERVER_NAME=

export TEMPORAL_INSECURE_SKIP_VERIFY=false

# payload data encryption
export ENCRYPT_PAYLOAD=false
export DATACONVERTER_ENCRYPTION_KEY_ID=mysecretkey

# App temporal taskqueue name
export TRANSFER_MONEY_TASK_QUEUE="go-moneytransfer"
export BANK_SERVICE_AVAILABLE="true"

# Set to enable debug logger logging
export LOG_LEVEL=debug

# local mysql backend db connection
export MYSQL_HOST=localhost
export MYSQL_DATABASE=dataentry
export MYSQL_USER=mysqluser
export MYSQL_PASSWORD=mysqlpw
```

Note: If you simulate a banking service outage on a deposit activity the BANK_SERVICE_AVAILABLE env is read by the worker on startup so changing the env required a worker restart.

## Start the webapp and navigate to view the local tasks

Start the webapp, by default it listens on port localhost:8085
```
go run webapp.go
```

Note: The webapp has a background activity thread that periodically polls the database transfer table on a go cron timer looking for entries with status "REQUESTED", if an entry is found it reads the oldest and updates it to "PROCESSING" and starts a corresponding Temporal Workflow.  Example polling can be seen in the webapp terminal:

```
go run webapp.go
2023/06/06 11:36:22 Serve Http on 8085
2023/06/06 11:36:52 CheckTransferQueueTask: called
2023/06/06 11:36:52 QueryTransferRequest: called
2023/06/06 11:36:52 CheckTransferQueueTask: No transfers in queue.
2023/06/06 11:37:22 CheckTransferQueueTask: called
2023/06/06 11:37:22 QueryTransferRequest: called
2023/06/06 11:37:22 CheckTransferQueueTask: No transfers in queue.
...
```

When the transfer poll task does find a pending entry the log looks like:
```
2023/06/06 11:39:03 CheckTransferQueueTask: called
2023/06/06 11:39:03 QueryTransferRequest: called
2023/06/06 11:39:04 QueryTransferRequest: Transfer: {4 bill ted 10 transfer request PROCESSING}
2023/06/06 11:39:04 CheckTransferQueueTask: Transfer: {4 bill ted 10 transfer request PROCESSING}
2023/06/06 11:39:04 CheckTransferQueueTask: PaymentDetails: {bill ted transfer request 10}
2023/06/06 11:39:04 StartMoneyTransfer-72937: called, PaymentDetails: moneytransfer.PaymentDetails{SourceAccount:"bill", TargetAccount:"ted", ReferenceID:"transfer request", Amount:10}
2023/06/06 11:39:04 StartMoneyTransfer-72937: Starting moneytransfer workflow on go-moneytransfer task queue
2023/06/06 11:39:04 StartMoneyTransfer-72937: Started workflow: WorkflowID: go-txfr-webtask-wkfl-72937, RunID: 09d781d2-8d95-4926-a458-b6f2ac49ad37
2023/06/06 11:39:33 CheckTransferQueueTask: called
2023/06/06 11:39:33 QueryTransferRequest: called
2023/06/06 11:39:34 CheckTransferQueueTask: No transfers in queue.
2023/06/06 11:40:03 CheckTransferQueueTask: called
2023/06/06 11:40:03 QueryTransferRequest: called
2023/06/06 11:40:04 CheckTransferQueueTask: No transfers in queue.
2023/06/06 11:40:09 StartMoneyTransfer-72937: Workflow result: "Transfer complete (transaction IDs: W4081768757, D2204353735)"
2023/06/06 11:40:09 StartMoneyTransfer-72937: done.
2023/06/06 11:40:09 UpdateTransferRequest: called (Id: 4 COMPLETED )
2023/06/06 11:40:09 CheckTransferQueueTask: Workflow: go-txfr-webtask-wkfl-72937 Completed
```

## Start the temporal Transfer worker

In a different terminal window to the webapp to separate the terminal log output displayed.

```
cd moneytransfer
go run worker/main.go

2023/06/06 11:27:05 Go worker starting..
2023/06/06 11:27:05 LoadClientOption: myns.abcdf.tmprl.cloud:7233 myns.abcdf ~/.temporal/tclient-myns.pem ~/.temporal/tclient-myns.key   false true
2023/06/06 11:27:05 Go worker connecting to server..
2023/06/06 11:27:05 Go worker initialising..
2023/06/06 11:27:05 Go worker registering for Workflow moneytransfer.Transfer..
2023/06/06 11:27:05 Go worker registering for Activity moneytransfer.Withdraw..
2023/06/06 11:27:05 Go worker registering for Activity moneytransfer.Deposit..
2023/06/06 11:27:05 Go worker registering for Activity moneytransfer.Refund..
2023/06/06 11:27:05 Go worker listening on go-moneytransfer task queue..
2023/06/06 11:27:05 INFO  Started Worker Namespace myns.abcdf TaskQueue go-moneytransfer WorkerID 94006@gmini.local@
```

Note: the sample worker registers with the taskqueue to handle Workflow and Activity actions for the Transfer Workflow example



