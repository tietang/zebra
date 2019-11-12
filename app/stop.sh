#!/usr/bin/env bash

ps -ef|grep zebra|awk '{print $2}'|while read pid
        do
                kill -9 $pid
        done