#!/bin/bash

# build
docker build -t local-mkdocs .

# run
exec docker run --rm -it -p 8000:8000 -v ${PWD}:/docs local-mkdocs $@
