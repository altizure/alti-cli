### Install or Update
```bash
go get -u github.com/jackytck/alti-cli

# auto-complete
alti-cli help completion
```

### Login
```bash
# email login
alti-cli login

# phone login (public api only)
alti-cli login -p
```
* Support public api-server, Altizure One and private api-server.
* e.g. endpoint for private server: http://1.2.3.4:1234

### Network Test
```bash
alti-cli network
```

### List Account
```bash
alti-cli account
```

### Switch Account
```bash
alti-cli account use XXXXXX
```
* XXXXXX is the account ID

### New Project
```bash
alti-cli project new recon -n 'test new proj'

+--------------------------+---------------+--------------+------------+
|            ID            |     NAME      | PROJECT TYPE | VISIBILITY |
+--------------------------+---------------+--------------+------------+
| 5d37e018bb7c6a0e17ffe9d1 | test new proj | free         | public     |
+--------------------------+---------------+--------------+------------+
```
* recon: reconstruction project
* -n: project name, e.g. 'test new proj'
* -s: output project id only

### Check local images
```bash
alti-cli check image -d ~/myimg -v -t -s .small -n 10
```
* -d: image directory, e.g. ~/myimg
* -v: verbose
* -t: table format
* -s: directory to skip, e.g. .small
* -n: number of threads, default is number of cores

### Import Image
```bash
alti-cli import image -d ~/myimg -s .small -p 5d37e -r upload.csv -v -m s3 -y
```
* -d: image directory, e.g. ~/myimg
* -s: directory to skip, e.g. .small
* -p: (partial) project id from aboved, e.g. 5d37e
* -r: name of report, e.g. upload.csv (not required)
* -v: verbose
* -m: upload method (skip this flag to auto detect best method)
* -n: number of threads, default is number of cores
* -y: auto accept

### Inspect Project
```bash
alti-cli myproj inspect -p 5d37e0

alti-cli myproj
```

### Start Reconstruction
```bash
alti-cli project start -p 5d37e0
```

### Download Results (pro project only)
```bash
alti-cli project download -p 5d37e -y
```

### Arbitrary GQL (query + mutation)
```bash
$ cat q.txt
query ($id: ID!) {
  project(id: $id) {
    id
    name
    isImported
    importedState
    projectType
    numImage
    gigaPixel
    taskState
    date
    cloudPath {
      key
    }
  }
}

$ cat var.txt
{
  "id": "5d37e018bb7c6a0e17ffe9d1"
}

$ alti-cli gql -q q.txt -k var.txt
{
  "project": {
    "cloudPath": [
      {
        "key": "local"
      },
      {
        "key": "oss_sz"
      },
      {
        "key": "s3"
      }
    ],
    "date": "2019-07-24T08:33:01.321Z",
    "gigaPixel": 0.34,
    "id": "5d37e018bb7c6a0e17ffe9d1",
    "importedState": "Native",
    "isImported": false,
    "name": "test new proj",
    "numImage": 28,
    "projectType": "free",
    "taskState": "Done"
  }
}
```
* -q: path of query or mutation file
* -k: path of query variables file
