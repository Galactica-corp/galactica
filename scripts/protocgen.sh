#!/usr/bin/env bash

# Copyright 2025 Galactica Network
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# How to run manually:
# docker build --pull --rm -f "contrib/devtools/Dockerfile" -t cosmossdk-proto:latest "contrib/devtools"
# docker run --rm -v $(pwd):/workspace --workdir /workspace cosmossdk-proto sh ./scripts/protocgen.sh

echo "Formatting protobuf files"
find ./ -name "*.proto" -exec clang-format -i {} \;

set -e

home=$PWD

echo "Generating proto code"
proto_dirs=$(find ./ -name buf.yaml -print0 | xargs -0 -n1 dirname | sort | uniq)
for dir in $proto_dirs; do
  echo "Generating proto code for $dir"
  
  cd "$dir"
  # check if buf.gen.pulsar.yaml exists in the proto directory
  if [ -f "buf.gen.pulsar.yaml" ]; then
    buf generate --template buf.gen.pulsar.yaml
    
    # move generated files to the right places
    if [ -d "../galactica" ] && [ "$dir" != "./proto" ]; then
      cp -r ../galactica "$home"/api
      rm -rf ../galactica
    fi
  fi
  
  # check if buf.gen.gogo.yaml exists in the proto directory
  if [ -f "buf.gen.gogo.yaml" ]; then
    buf generate --template buf.gen.gogo.yaml
    
    # move generated files to the right places
    if [ -d "../galactica" ]; then
      cp -r ../galactica/* "$home"
      rm -rf ../galactica
    fi
    
    if [ -d "../github.com" ] && [ "$dir" != "./proto" ]; then
      cp -r ../github.com/Galactica-corp/galactica/* "$home"
      rm -rf ../github.com
    fi
  fi
  
  cd "$home"
done

# move generated files to the right places
cp -r github.com/Galactica-corp/galactica/* ./
rm -rf github.com

go mod tidy
