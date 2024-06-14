#!/bin/bash

source scripts/shared.sh

# run
exec docker run \
    --rm \
    --name "${DOCKER_CONTAINER_NAME}" \
    --publish 8000:8000 \
    --volume "${PWD}:/docs" \
    "${DOCKER_IMAGE}" \
    $@
