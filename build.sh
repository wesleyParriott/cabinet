go build

if [ "$?" -eq 0 ]; then
    if [ "$1" = "run" ]; then
        ./cabinet
    fi
fi
