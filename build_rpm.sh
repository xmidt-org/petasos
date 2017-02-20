#!/bin/bash

echo "Adjusting build number..."

OIFS=$IFS
IFS='

'

release=""

taglist=`git tag -l`
tags=($taglist)

for ((i=${#tags[@]}-1; i >=0; i--)); do
    if [[ "${tags[i]}" != *"alpha"* ]]; then
        release=${tags[i]}
        break
    fi
done

if [ -z "$release"  ]; then
    echo "Could not find latest release tag!"
else
    echo "Most recent release tag: $release"
fi

IFS=$OIFS

new_release=`echo "$release" | awk -F. '{$NF+=1; OFS="."; print $0}'`
new_release+="-${BUILD_NUMBER}alpha"
echo "Issuing release $new_release..."

echo "Building the petasos rpm..."
docker exec build bash -c "pushd petasos; git fetch; git checkout travis-testing; popd"
docker exec build bash -c "mv petasos petasos-$new_release; tar -czvf petasos-$new_release.tar.gz petasos-$new_release; mv petasos-$new_release.tar.gz /root/rpmbuild/SOURCES"
docker exec build bash -c "pushd /root/rpmbuild; ls -R; popd"
docker exec build bash -c "pushd petasos-$new_release; rpmbuild -ba petasos.spec; popd"

