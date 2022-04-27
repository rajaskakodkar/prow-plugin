#!/bin/sh

# Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
# SPDX-License-Identifier: Apache-2.0

# set this value to your package name
NAME=$1
VERSION=$2

if [ -z "$NAME" ]
then
  echo "create package failed. must set NAME"
  exit 2
fi

if [ -z "$VERSION" ]
then
  echo "create package failed. must set VERSION"
  exit 2
fi

echo "\n===> Updating image lockfile for package ${NAME}/${VERSION}\n"
cd packages/${NAME}/${VERSION} && kbld --file bundle --imgpkg-lock-output bundle/.imgpkg/images.yml