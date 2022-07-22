#!/bin/bash

install_sqlite_then_import() {
    platform=$(uname -s)

    if [ $platform == "Darwin" ]
    then
        brew update ; brew upgrade
        brew install sqlite3
    elif [ $platform == 'Linux' ]
    then
        sudo apt update ; sudo apt install -y sqlite3
    else
        echo "Unsupport platform!"
    fi
}

main() {
    install_sqlite_then_import
}

main
