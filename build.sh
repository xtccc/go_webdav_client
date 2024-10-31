#!/bin/bash
linux(){
    rm -f go_webdav_client
    go build -o doc/go_webdav_client .
    if [ $? = 0 ] ;then
        ./go_webdav_client upload -f ./go_webdav_client -w /go_webdav_client_test
        sudo cp go_webdav_client /usr/local/bin/go_webdav_client
    else
        echo "build error"
    fi
}



# https://192.168.31.175/go_webdav_client_test/



android(){
    export GOOS=linux
    export  GOARCH=arm64
    go build -o doc/go_webdav_client_arm64 . && ./go_webdav_client upload -f ./go_webdav_client_arm64 -w /go_webdav_client_test
    unset GOOS
    unset GOARCH
}

function usage(){
    echo "all to upload all"
    echo "两个参数: linux android "
}
case $1 in
    linux)
    linux;;
    android)
    android;;
    all)
    linux && android ;;
    *)
    usage;;
esac
