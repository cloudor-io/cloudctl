#!/bin/bash
supported_os=("linux" "darwin" "windows")
suffix=("linux_x86_64" "macosx_x86_64" "windows_x86_64.exe")
supported_arch=("amd64")
release_root=/opt/cloudor/cmd/releases
for os_index in "${!supported_os[@]}"
do
    os=${supported_os[$os_index]}
    for arch in "${supported_arch[@]}"
    do
        mkdir -p $release_root/$os/$arch/latest/
        echo "build to $release_root/$os/$arch/latest/cloudor"
        CGO_ENABLED=0 GOOS=$os GOARCH=$arch go build -ldflags "-s -w" -o $release_root/$os/$arch/latest/cloudor main.go
        chmod +x $release_root/$os/$arch/latest/cloudor
        echo "copy to $release_root/$os/$arch/latest/cloudor_${suffix[$os_index]}"
        cp $release_root/$os/$arch/latest/cloudor $release_root/$os/$arch/latest/cloudor_${suffix[$os_index]}
    done        

done