#!/usr/bin/env bash
set -e

source scripts/shared.sh

project_path=$(git rev-parse --show-toplevel)

exec docker run \
    --rm \
    --name $DOCKER_CONTAINER_NAME \
    --volume "${project_path}:/project" \
    --workdir /project/scripts/servers \
    --entrypoint '' \
    $DOCKER_IMAGE \
    /bin/masto-guide-dk-servers build
