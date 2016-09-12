# Cloud-Config-Server

A service to serve your [cloud-init](https://cloudinit.readthedocs.io/en/latest/) configurations using golang templates.

You can start the service as follows:

```
docker run -d -p 8080:8080 -e "PORT=8080" -e "WORKDIR=/cloud-init" -v $PWD:/cloud-init ictu/cloud-config-server
```

The ```WORKDIR``` variable specifies where the cloud init template files reside.

From your favourite http client ```GET```
```
curl http://localhost:8080/cloud-init/example
```

will parse the ```example.yml``` template using ```example-vars.yml``` as data object.

Add as many templates as you like!
