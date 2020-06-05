#!/usr/bin/env bash

# Script useful for manage stack services life-cycle.
# See usage for commands & options.

set -o errexit;

readonly STACK_NAME=workflows
readonly DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
readonly VERSION="$(cat < ./.env | grep VERSION_TAG= | cut -d= -f2)"
readonly DOCKER_HUB_REPO=emanueljoivo/workflows
readonly POSTGRES_DB=workflows-db
readonly POSTGRES_USER=workflows-admin
readonly POSTGRES_PASSWORD="$(cat < ./.env | grep DATABASE_PASSWORD= | cut -d= -f2)"
readonly POSTGRES_PORT=5432

usage ()
{
    printf "usage: %s run | release | build | deploy | publish | clean | rm | [-h]\n" "$0"
}

service_build()
{
  local DOCKERFILE_DIR="$DIR"/Dockerfile

  docker build -t "${DOCKER_HUB_REPO}":"${VERSION}" \
            --file "${DOCKERFILE_DIR}" .
}

db_run() {
  docker run -it -p "$POSTGRES_PORT":"$POSTGRES_PORT" -v "$POSTGRES_DB":/var/lib/postgresql/data \
                -e POSTGRES_USER="$POSTGRES_USER"\
                -e POSTGRES_PASSWORD="$POSTGRES_PASSWORD"\
                -e POSTGRES_DB="$POSTGRES_DB" postgres
}

workflows_run() {
  docker run -it -p 5000:5000 --name workflows emanueljoivo/workflows:"$VERSION"
}

service_run() {
  case $1 in
    db) shift
      db_run
      ;;
    workflows) shift
      workflows_run
      ;;
    *)
      echo "No way to run this"
      ;;
  esac
}

stack_deploy()
{
  docker stack deploy -c "${DIR}/"docker-stack.yml "${STACK_NAME}"
}

stack_publish()
{
  docker push "${DOCKER_HUB_REPO}":"${VERSION}"
}

stack_clean()
{ 
  docker stack rm "$STACK_NAME"
  docker system prune -f
}

stack_rm()
{
  docker stack rm "${STACK_NAME}"
}

release() {
  local VERSION_TAG=$1
  local VERSION_NAME=$2
  sed -i "/^VERSION_TAG/c\VERSION_TAG=${VERSION_TAG}" ./.env
  sed -i "/^VERSION_NAME/c\VERSION_NAME=${VERSION_NAME}" ./.env
  service_build
  stack_publish
}

define_params()
{
    case $1 in
        run) shift
          service_run "$@"
          ;;
        release) shift
            release "$@"
            ;;
        build) shift
            service_build
            ;;
        deploy) shift
            stack_deploy
            ;;
        publish) shift
            service_build
            stack_publish
            ;;
        clean) shift
            stack_clean
            ;;
        rm) shift
            stack_rm
            ;;
        -h | --help | *)
            usage;
            exit 0;
            ;;
    esac
}

define_params "$@"
