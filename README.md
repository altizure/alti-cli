### Email Login
```bash
./alti-cli login
```

### Phone Login
```bash
./alti-cli login -p
```

### List Account
```bash
./alti-cli account
```

### Switch Account
```bash
./alti-cli account use XXXXXX
```
* XXXXXX is the account ID

### New Project
```bash
./alti-cli project new recon -n 'test new proj'

+--------------------------+---------------+--------------+------------+
|            ID            |     NAME      | PROJECT TYPE | VISIBILITY |
+--------------------------+---------------+--------------+------------+
| 5d37e018bb7c6a0e17ffe9d1 | test new proj | free         | public     |
+--------------------------+---------------+--------------+------------+
```
* recon: reconstruction project
* -n: project name, e.g. 'test new proj'
* -s: to output project id only

### Import Image
```bash
./alti-cli import image -d ~/myimg -s .small -p 5d37e -r upload.csv -v -m s3 -y
```
* -d: image directory, e.g. ~/myimg
* -s: directory to skip, e.g. .small
* -p: (partial) project id from aboved, e.g. 5d37e
* -r: name of report, e.g. upload.csv (not required)
* -v: is verbose
* -m: is upload method (skip this flag to auto detect best method)
* -y: auto accept

### Inspect
```bash
./alti-cli myproj inspect -p 5d37e0

./alti-cli myproj
```

### Start Reconstruction
```bash
./alti-cli project start -p 5d37e0
```
