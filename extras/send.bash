#!/bin/bash
for i in {1..50}
do
  ./sampleRabbitMqSendMessage -message="bom pesan ke $i dari 50 dan jede 1 menit" -target="6282329949292-1590306644@g.us"
  sleep 1m
done
