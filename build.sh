image_version=$1

echo "Compile the application"
docker run --rm -v $(pwd):/usr/src/myapp -w /usr/src/myapp -e CGO_ENABLED=0 -e GOOS=linux -e GOARCH=amd64 golang:1.7.1 bash -c "go get -d -v; go build -a --installsuffix cgo -v -o cloud-config-server"

echo "Building cloud-config-server:$image_version"
docker build --no-cache=true -t ictu/cloud-config-server:$image_version .

echo "Pushing cloud-config-server:$image_version"
docker push ictu/cloud-config-server:$image_version

echo "Cleanup"
rm -f cloud-config-server
