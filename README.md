[![Go Report Card](https://goreportcard.com/badge/github.com/jackytck/alti-cli)](https://goreportcard.com/report/github.com/jackytck/alti-cli)

### Install or Update
```bash
$ go get -u github.com/jackytck/alti-cli

# bash auto-complete
$ alti-cli help completion
```

### Login
```bash
# email login
$ alti-cli login

# phone login (public api only)
$ alti-cli login -p

# login with specific key (e.g. your paid developer key)
$ alti-cli login -k
```
* Support public api-server, Altizure One and private api-server.
* e.g. endpoint for private server: http://1.2.3.4:1234

### Quick start
1. Put all images and meta files in a directory (e.g. /tmp/ust-test), or zipped obj (e.g. /tmp/bunny.zip)
2. Call
```bash
# create project + upload (images/metafiles/model) + start task
$ alti-cli quick -i /tmp/ust-test

# Or a zip file of CAD obj
$ alti-cli quick -i /tmp/bunny.zip
```
3. Done

### Network Test
Check if direct upload is supported.
```bash
$ alti-cli network
```

### Site Test
Check if main browsing site is up.
```bash
$ alti-cli check site -v

2019/11/06 14:56:45 Checking https://www.altizure.com...
2019/11/06 14:56:45 Success with status code: 200
2019/11/06 14:56:45 Took 637.822755ms
```

### List Account
```bash
$ alti-cli account
```

### Switch Account
```bash
$ alti-cli account use XXXXXX
```
* XXXXXX is the account ID

### Current active user
```bash
$ alti-cli whoami

+--------------------------+----------------+--------+----------+------------+---------------+--------------+------+---------+--------+------+-----------+---------------------+---------------+
|         ENDPOINT         | USERNAME/EMAIL | COINS  | GP QUOTA | MEMBERSHIP | DEVELOPERSHIP | FREE STORAGE | STAR | PROJECT | PLANET | FANS | FOLLOWING |       JOINED        |    COUNTRY    |
+--------------------------+----------------+--------+----------+------------+---------------+--------------+------+---------+--------+------+-----------+---------------------+---------------+
| https://api.altizure.com | jacky          | 105.32 |     9.00 | ACTIVE     | ACTIVE        | 3.8 GiB      |    1 |     124 |      2 |    2 |         4 | 2015-08-21 09:33:07 | Hong_Kong_SAR |
+--------------------------+----------------+--------+----------+------------+---------------+--------------+------+---------+--------+------+-----------+---------------------+---------------+
```

### Mmebership info
```bash
$ alti-cli my membership

+--------+-------+--------+---------------------+---------------------+----------+---------+---------+---------+------------+-----------+---------------+--------------+-----------+
| STATE  | PLAN  | MONTHS |        START        |         END         | GP QUOTA | COIN/GP | STORAGE |  USAGE  | VISIBILITY |  COUPON   | MODEL/PROJECT | COLLABORATOR | WATERMARK |
+--------+-------+--------+---------------------+---------------------+----------+---------+---------+---------+------------+-----------+---------------+--------------+-----------+
| ACTIVE | Small |     12 | 2018-12-24 05:37:55 | 2019-12-24 05:37:55 |     6.00 |    1.00 | 4.0 GiB | 191 MiB | public     | Value: 5  |            10 |            5 | true      |
|        |       |        |                     |                     |          |         |         |         |            | Repeat: 1 |               |              |           |
|        |       |        |                     |                     |          |         |         |         |            | Month: 1  |               |              |           |
+--------+-------+--------+---------------------+---------------------+----------+---------+---------+---------+------------+-----------+---------------+--------------+-----------+
```

### New Project (reconstruction)
```bash
$ alti-cli project new recon -n 'test new proj'

+--------------------------+---------------+--------------+------------+
|            ID            |     NAME      | PROJECT TYPE | VISIBILITY |
+--------------------------+---------------+--------------+------------+
| 5d37e018bb7c6a0e17ffe9d1 | test new proj | free         | public     |
+--------------------------+---------------+--------------+------------+
```
* recon: reconstruction project
* -n: project name, e.g. 'test new proj'
* -s: output project id only

### New Project (imported model)
```bash
$ alti-cli project new model -n 'my obj model'
```

### Remove Project (any kind)
```bash
$ alti-cli project remove -p '5d37e018bb7c6a0e17ffe9d1'
```

### Check local images (without uploading)
Check all images of a given directory locally. Get stats of number of GP, dimensions and invalid images, etc.
```bash
$ alti-cli check image -d ~/myimg -v -t -s .small -n 10
```
* -d: image directory, e.g. ~/myimg
* -v: verbose
* -t: table format
* -s: directory to skip, e.g. .small
* -n: number of threads, default is number of cores

### List buckets
Buckets are used in the `import` command for specifying different geo endpoints for the upload process. Would be auto selected if not provided.
```bash
$ alti-cli list bucket
```

### Import Image (reconstruction project)
```bash
$ alti-cli import image -d ~/myimg -s .small -p 5d37e -r upload.csv -v -m s3 -y
```
* -d: image directory, e.g. ~/myimg
* -s: directory to skip, e.g. .small
* -p: (partial) project id from aboved, e.g. 5d37e
* -r: name of report, e.g. upload.csv (not required)
* -v: verbose
* -m: upload method (skip this flag to auto detect best method)
* -n: number of threads, default is number of cores
* -y: auto accept

### Import Meta file (reconstruction project)
```bash
$ alti-cli import meta -p 5d008 -v -f ~/test/pose.txt
```
* -b: desired bucket to upload (auto select if empty)
* -f: path of meta file
* -p: (partial) project id from aboved, e.g. 5d37e
* -m: method of upload: 'direct' or 's3' or 'minio' (based on supported cloud shown in `alti-cli account`)
* -t: timeout in second(s)
* -ip: ip address of ad-hoc local server for direct upload
* -port: port of ad-hoc local server for direct upload
* -v: verbose

### Import Model file (imported model project)
```bash
$ alti-cli import model -p 5d7b6b -v -f ~/test/bunny.obj
```
* -b: desired bucket to upload
* -f: path of model zip file or directory of multiparts zip
* -p: (partial) project id from aboved, e.g. 5d37e
* -m: method of upload: 'direct' or 's3' or 'minio'
* -t: timeout in second(s)
* -v: verbose

### Inspect Project
```bash
$ alti-cli myproj inspect -p 5d37e0

$ alti-cli myproj
```

### Start Reconstruction
```bash
$ alti-cli project start -p 5d37e0
```

* -t: task type: One of `alti-cli list task-type`, default is `Native`

### List task types
```bash
$ alti-cli list task-type
```

### Stop Reconstruction
```bash
$ alti-cli project stop -p 5d37e0
```

### Download Results (pro project only)
```bash
$ alti-cli project download -p 5d37e -y
```

### Bank
Cash to coins
```bash
$ alti-cli bank tocoins -c 12

USD12.00 could buy 1.02 coins
```

Coins to cash
```bash
$ alti-cli bank tocash -c 3

3.00 coins could be bought by USD35.40
```

P2P coins
```bash
$ alti-cli bank transfer -c 0.1 -e nat@nat.com -m 'Test transferring coins'

You have 105.42 coins. Are you sure to transfer 0.10 coins to "nat@nat.com"? (Y/N): y
Successfully transferred 0.10 coins to "nat@nat.com"
Current balnce: 105.32 coins
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

### Super tools

#### Get user token
```bash
$ alti-cli super token -e 'nat@altizure.com'
```
* -e: user email

#### Sync project to cloud
```bash
$ alti-cli super sync -p 5d7b6b
```
* -p: (partial) project id

#### Sync images from cloud to gfs
```bash
$ alti-cli super toGFS -p 5d7b6b
```
* -p: (partial) project id

#### Force start project task
```bash
$ alti-cli super force-start -p 5d7b6b
```
* -p: (partial) project id
