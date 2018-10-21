# gsa (grootfs store accounting? maybe?)
how much disk is my grootfs store using?

this is a toy which i built in an hour while waiting for a kernel to compile, it will probably fail :)

```
$ wget https://github.com/Callisto13/gsa/releases/download/v1.1/gsa -P /usr/local/bin && chmod +x /usr/local/bin/gsa

$ gsa --help
Usage of gsa:
  -grootfs-bin string
        path to the grootfs bin (default "/var/vcap/packages/grootfs/bin/grootfs")
  -grootfs-config string
        path to grootfs' config (default "/var/vcap/jobs/garden/config/grootfs_config.yml")
  -r    human readable result

$ gsa
{"total_bytes_containers":95608832,"total_bytes_layers":893982999,"total_bytes_active_layers":893982999,"total_bytes_store":989591831}

$ gsa -r
Containers: 96 MB
Layers: 894 MB (of which Active: 894 MB)
990 MB
```

note: grootfs only tracks the amount of container disk usage if the rootfs was created with a disk limit. if your containers (and rootfses) were made without limits, then grootfs will always report `"total_bytes_containers":0`.
