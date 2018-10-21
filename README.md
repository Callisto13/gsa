# gsa (grootfs store accounting? maybe?)
how much disk is my grootfs store using?

this is a toy which i built in an hour one sunday morning, it will probably fail :)

```
$ wget https://github.com/Callisto13/gsa/releases/download/v1.0/gsa -P /usr/local/bin && chmod +x /usr/local/bin/gsa
$ gsa --help
Usage of gsa:
  -grootfs-bin string
        path to the grootfs bin (default "/var/vcap/packages/grootfs/bin/grootfs")
  -grootfs-config string
        path to grootfs' config (default "/var/vcap/jobs/garden/config/grootfs_config.yml")
$ gsa
{"total_bytes_containers":95608832,"total_bytes_layers":893982999,"total_bytes_active_layers":893982999}
```
