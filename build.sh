#!/bin/sh

if [ ! -f admin ]; then
	rm admin
fi
cd cmd/admin
go build -o admin
cd ../../
mv cmd/admin/admin .
chmod +x admin
./admin
