#!/usr/bin/env bash
SCRIPT=$(readlink -f $0)
BASEDIR=`dirname $SCRIPT`
docker run -v $BASEDIR/data:/data/data -v $BASEDIR/config.yaml:/data/config.yaml -v /var/run/docker.sock:/var/run/docker.sock parsermanager ./app