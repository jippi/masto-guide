export DOCKER_IMAGE="ghcr.io/jippi/masto-guide:main"
export DOCKER_CONTAINER_NAME=masto-guide-local

if [ -z $CI ]
then
    export DOCKER_IMAGE=$DOCKER_CONTAINER_NAME

    # build
    docker build -t $DOCKER_IMAGE .
fi
