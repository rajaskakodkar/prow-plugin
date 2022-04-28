#!/bin/bash

OCI_REGISTRY="public.ecr.aws/t0q8k6g2/prow-plugin"
package="$1"
version="$2"

imgpkg push --bundle ${OCI_REGISTRY}/${package}:${version} --file packages/${package}/${version}/bundle

