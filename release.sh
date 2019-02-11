#!/usr/bin/env bash

package_name=ajz_xyz/experimental/computation/mfm-go/mfmgo

platforms=(
"windows/386"
"windows/amd64"
"darwin/386"
"darwin/amd64"
"linux/386"
"linux/amd64"
)

mkdir -p bin

for platform in "${platforms[@]}"
do
    platform_split=(${platform//\// })
    GOOS=${platform_split[0]}
    GOARCH=${platform_split[1]}
    output_name='mfmgo_'$GOOS'_'$GOARCH
    if [ $GOOS = "windows" ]; then
        output_name+='.exe'
    fi  

    GOOS=$GOOS GOARCH=$GOARCH go build -o bin/$output_name $package_name
    if [ $? -ne 0 ]; then
        echo 'An error has occurred! Aborting the script execution...'
        exit 1
    fi
done
