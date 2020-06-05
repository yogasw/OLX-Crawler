#!/bin/bash

echo -n "WARNING! delete all docker images and containers on your machine (No/yes) : "
read ANSWER

if [[ $ANSWER == "yes" ]]
then
  # Delete every Docker containers
  # Must be run first because images are attached to containers

  docker rm -f $(docker ps -a -q)

  # Delete every Docker image
  docker rmi -f $(docker images -q)

#  sudo rm -rf .data/
fi
