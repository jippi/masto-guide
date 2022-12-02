#!/bin/bash

DOCKER_IMAGE="ghcr.io/jippi/masto-guide:main"

if [ -z $CI ]
then
    DOCKER_IMAGE=local-mkdocs

    # build
    docker build -t local-mkdocs .
fi

# run
exec docker run --name mkdocs-local --rm -p 8000:8000 -v ${PWD}:/docs $DOCKER_IMAGE $@
