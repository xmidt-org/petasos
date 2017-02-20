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

echo "Most recent release tag: $release"

IFS=$OIFS

# # begin old script
# echo "Hello world."
# touch versionno.txt
# 
# echo $(git tag -l)
# 
# git for-each-ref --sort=-authordate refs/tags | \
#     while read entry; do
#         echo $entry
#         tag=`echo $entry | awk '{print $NF}'`
#         tag=`echo $tag | awk -F/ '{print $NF}'`
#         if [[ "$tag" != *"-"* ]]; then
#             echo "$tag" > versionno.txt
#             break
#         fi
#     done
# 
# release=`cat versionno.txt`
# rm versionno.txt
# 
# filename=`echo "$release" | awk -F. '{$NF+=1; OFS="."; print $0}'`
# filename+="-${BUILD_NUMBER}alpha"
# echo $filename
# 
# # echo "Building the petasos rpm..."
# # docker exec build bash -c "pushd petasos; git fetch; git checkout travis-testing; popd"
# # docker exec build bash -c "mv petasos petasos-1.1.0; tar -czvf petasos-1.1.0.tar.gz petasos-1.1.0; mv petasos-1.1.0.tar.gz /root/rpmbuild/SOURCES"
# # docker exec build bash -c "pushd petasos-1.1.0; rpmbuild -ba petasos.spec; popd"
# 
