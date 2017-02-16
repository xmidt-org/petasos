#!/bin/bash

echo "Hello world."

release=`git describe --abbrev=0 --tags`

filename=`echo "$release" | awk -F. '{$NF+=1; OFS="."; print $0}'`
filename+="-${BUILD_NUMBER}alpha"

echo $release
echo $filename

# echo "Building the petasos rpm..."
# docker exec build bash -c "pushd petasos; git fetch; git checkout travis-testing; popd"
# docker exec build bash -c "mv petasos petasos-1.1.0; tar -czvf petasos-1.1.0.tar.gz petasos-1.1.0; mv petasos-1.1.0.tar.gz /root/rpmbuild/SOURCES"
# docker exec build bash -c "pushd petasos-1.1.0; rpmbuild -ba petasos.spec; popd"

