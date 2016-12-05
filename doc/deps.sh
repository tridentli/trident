#!/bin/sh -e

echo "== Trident dependency updater"

TRIDENT=$(pwd)

if [ ! -f "${TRIDENT}/doc/deps.sh" ];
then
	echo "Please run from root of Trident as doc/deps.sh"
	exit 1
fi

# Our Go Path
mkdir -p ${GOPATH}

echo "Trident path: ${TRIDENT}"
echo "Trident GOPATH: ${GOPATH}"

# Update or Fetch EpicEditor
echo "- EpicEditor"
if [ -d ${TRIDENT}/ext/epiceditor ];
then
	cd ${TRIDENT}/ext/epiceditor/
	git pull
else
	git clone https://github.com/OscarGodson/EpicEditor.git ${TRIDENT}/ext/epiceditor
fi

# Fetch & Update go dependencies
echo "- Go Dependencies"

cd ${TRIDENT}
go get -d ./...
go get -d -u all || true

echo "== Trident dependency updater -- done"

