#!/usr/bin/env bash
set -e

project_path=$(git rev-parse --show-toplevel)
docker run --rm -v "${project_path}:/project" -w /project/scripts/servers golang:1.19-alpine go run .
