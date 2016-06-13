#!/bin/sh
mkdir -p logs
mkdir -p pids

rm  pids/*

# DB
go install github.com/ezeev/golang/traderdash/dbi

# Yahoo Finance Loader
echo Installing and Running Yahoo Finance Quote Loader
go install github.com/ezeev/golang/traderdash/finance
$GOPATH/bin/finance -symbols AAPL,MSFT,GOOGL -interval 10 -db-type rethinkdb -conn-str localhost:28015 > logs/finance.log 2>&1 &
FIN_PID=$!
echo Started Finance PID: $FIN_PID
echo $FIN_PID >> pids/finance.pid

# Twitter Subscriber
echo Installing and Running Twitter Subscriber
go install github.com/ezeev/golang/traderdash/twitter-sub
$GOPATH/bin/twitter-sub -nats-urls http://192.168.99.100:4222 -db-type rethinkdb -conn-str localhost:28015 > logs/twitter-sub.log 2>&1 &
TWSUB_PID=$!
echo Started Twitter Sub PID: $TWSUB_PID
echo $TWSUB_PID >> pids/twitter-sub.pid

# Twitter Publisher
echo Installing and Running Twitter Publisher
go install github.com/ezeev/golang/traderdash/twitter-pub
$GOPATH/bin/twitter-pub -consumer-key W6OWcUnzlM2OkHmKwKLip9BCV \
 -consumer-secret XiYyh7TJZ0n01ZTdhiRJb6Fa0XxAuIGAfc8DeGfrwZx0Falb8A \
 -access-token 2944629920-hHZ8ADHX5S2ZamZJObwWFkijirabCxrWO2AIZAi \
-access-secret  9TMtULc2vEpFguK0ml2jXl1bao2XbObGfraUTsQXG4pYU \
-tracks donald,hillary \
-nats-urls http://192.168.99.100:4222 > logs/twitter-pub.log 2>&1 &
TWPUB_PID=$!
echo Started Twitter Sub PID: $TWPUB_PID
echo $TWPUB_PID >> pids/twitter-pub.pid
