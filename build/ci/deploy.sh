#!/usr/bin/env bash

set -euo pipefail

curl -sL https://git.io/goreleaser | bash
echo $DOCKER_PASSWORD | docker login -u $DOCKER_USERNAME --password-stdin
make push VERSION=${TRAVIS_BRANCH}
