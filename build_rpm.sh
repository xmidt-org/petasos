#!/bin/bash

echo "Hello world."
tag="temp"

git for-each-ref --sort=-authordate refs/tags | \
    while read entry; do
        tag=`echo $entry | awk '{print $NF}'`
        tag=`echo $tag | awk -F/ '{print $NF}'`
        if [[ "$tag" != *"-"* ]]; then
            echo "Found previous release: $tag..."
            break
        fi
    done

echo "$tag"
filename=`echo "$tag" | awk -F. '{$NF+=1; OFS="."; print $0}'`
echo $filename
filename+="-${BUILD_NUMBER}alpha"
echo $filename

# echo "Building the petasos rpm..."
# docker exec build bash -c "pushd petasos; git fetch; git checkout travis-testing; popd"
# docker exec build bash -c "mv petasos petasos-1.1.0; tar -czvf petasos-1.1.0.tar.gz petasos-1.1.0; mv petasos-1.1.0.tar.gz /root/rpmbuild/SOURCES"
# docker exec build bash -c "pushd petasos-1.1.0; rpmbuild -ba petasos.spec; popd"

