#!/bin/bash

EXECUTABLE_PATH="/home/shane/Dev/GO/go-pckg-manager/go-pckg-manager"

execute_command(){
    command="$1"
    shift

    "$EXECUTABLE_PATH" "$command" "$@"

}

if [ ! -x "$EXECUTABLE_PATH" ]; then
    echo "Go executable not found at $EXECUTABLE_PATH"
    exit 1
fi

command="$1"
shift

case "$command" in
    install)
	execute_command "install" "$@"
	;;
    remove)
	execute_command "remove" "$@"
	;;
    update)
	execute_command "update" "$@"
	;;
    search)
	execute_command "search" "$@"
	;;
    list)
	execute_command "list" "$@"
	;;
    *)
	echo "Invalid command: $command"
	echo "Usage: $0 [command] [arguments]"
	exit 1
	;;
esac
