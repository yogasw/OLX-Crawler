#!/bin/bash
for i in {1..50}
do
  ./sampleRabbitMqSendMessage -message="bom pesan ke $i dari 50 dan jede 1 menit" -target="6282127618761@s.whatsapp.net"
  sleep 1m
done
