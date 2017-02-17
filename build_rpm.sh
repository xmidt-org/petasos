#!/bin/bash

echo "Hello world."
touch versionno.txt

git for-each-ref --sort=-authordate refs/tags | \
    while read entry; do
        tag=`echo $entry | awk '{print $NF}'`
        tag=`echo $tag | awk -F/ '{print $NF}'`
        if [[ "$tag" != *"-"* ]]; then
            echo "$tag" > versionno.txt
            break
        fi
    done

release=`cat versionno.txt`
rm versionno.txt

filename=`echo "$release" | awk -F. '{$NF+=1; OFS="."; print $0}'`
filename+="-${BUILD_NUMBER}alpha"
echo $filename

# echo "Building the petasos rpm..."
# docker exec build bash -c "pushd petasos; git fetch; git checkout travis-testing; popd"
# docker exec build bash -c "mv petasos petasos-1.1.0; tar -czvf petasos-1.1.0.tar.gz petasos-1.1.0; mv petasos-1.1.0.tar.gz /root/rpmbuild/SOURCES"
# docker exec build bash -c "pushd petasos-1.1.0; rpmbuild -ba petasos.spec; popd"

