go build

if [ "$?" -eq 0 ]; then
    if [ "$1" = "run" ]; then
        ./cabinet
    fi
    if [ "$1" = "package" ]; then
        mkdir -p cabinet_package
        mkdir -p cabinet_package/setup
        mkdir -p cabinet_package/tmpls
        cp ./cabinet ./cabinet_package/
        cp ./setup/* ./cabinet_package/setup/
        cp ./tmpls/* ./cabinet_package/tmpls/
        tar zcf cabinet.tar.gz cabinet_package/
        rm -r cabinet_package
    fi
fi
