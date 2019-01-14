#!/bin/sh

LEVEL=2
# retrieve branch name
BRANCH_NAME=$(git branch | sed -n '/\* /s///p')

echo "Current version branch is $BRANCH_NAME"

# retrieve the last commit on the branch
VERSION=$(git describe --tags --match=$BRANCH_NAME* --abbrev=0)

if [ -z "$VERSION" ];then
    VERSION="0.0.0"
fi

NEW_TAG=
div=
i=0
for SUB in $(echo $VERSION | tr "\." "\n")
do
	if [[ $i == $LEVEL ]]
	then
		NEW_TAG=$NEW_TAG$div$(expr $SUB + 1)
		i=$(expr $i + 1)
		break
	fi
	
	NEW_TAG=$NEW_TAG$div$SUB
	div=.
	i=$(expr $i + 1)
done

echo "Updating $VERSION to $NEW_TAG"

# #get current hash and see if it already has a tag
GIT_COMMIT=`git rev-parse HEAD`
NEEDS_TAG=`git describe --contains $GIT_COMMIT`

#only tag if no tag already (would be better if the git describe command above could have a silent option)
if [ -z "$NEEDS_TAG" ]; then
    echo "Tagged with $NEW_TAG (Ignoring fatal:cannot describe - this means commit is untagged) "
    export BUILD_TAG=$NEW_TAG
    git tag -m "Bump version to $NEW_TAG" $NEW_TAG
    git push --tags
else
    echo "Already a tag on this commit"
fi
