### Login
```bash
./alti-cli login
```


### New Project
```bash
./alti-cli project new recon -n 'test new proj'

+--------------------------+---------------+--------------+------------+
|            ID            |     NAME      | PROJECT TYPE | VISIBILITY |
+--------------------------+---------------+--------------+------------+
| 5d37e018bb7c6a0e17ffe9d1 | test new proj | free         | public     |
+--------------------------+---------------+--------------+------------+
```

### Import Image
```bash
./alti-cli import image -d ~/myimg -s .small -p 5d37e -r upload.csv -v -m s3
```
* myimg: image directory
* .small: skip the .small directory
* 5d37e is the pid returned from above
* upload.csv is the upload report
* -v is verbose
* -m is upload method

### Inspect
```bash
./alti-cli myproj inspect -p 5d37e0

./alti-cli myproj
```
