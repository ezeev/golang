#!/bin/sh
echo Stopping finance
kill -9 `cat pids/finance.pid`
echo Stopping twitter pub
kill -9 `cat pids/twitter-pub.pid`
echo Stopping twitter sub
kill -9 `cat pids/twitter-sub.pid`
