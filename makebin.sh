#!/bin/bash
# NOTE : Quote it else use array to avoid problems #
FILES="./*.go"
/bin/rm -rf ./bin/*
for f in $FILES
do
  filename="${f%.*}"
  echo "Processing $filename file..."

  # macos
  GOOS=darwin GOARCH=amd64 go build -o ./bin/$filename-macos $filename.go

  # windows 
  GOOS=windows GOARCH=amd64 go build -o ./bin/$filename-64.exe $filename.go
  GOOS=windows GOARCH=386 go build -o ./bin/$filename-32.exe $filename.go

  # linux 
  GOOS=linux GOARCH=amd64 go build -o ./bin/$filename-linux $filename.go
done

cd ./bin
/usr/bin/zip CT-windows32.zip *-32.exe
/usr/bin/zip CT-windows64.zip *-64.exe
/usr/bin/zip CT-linux *-linux
/usr/bin/zip CT-macos *-macos
/bin/rm *.exe
/bin/rm *-linux
/bin/rm *-macos

echo ""
echo "Zip files are ready in ./bin"
