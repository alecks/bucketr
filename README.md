# bucketr

**bucketr**est is a lightweight RESTful API service that exposes contained key-value databases to the public. As of now, bucketr supports the following databases:

- [x] MongoDB
- [ ] etcd
- [ ] BoltDB (bbolt)
- [ ] Redis

It's still in early development, so database options are limited. DBs aren't particularly listed in order of priority, but etcd is the most demanded.

## Aims

- Make getting/setting values the exact same API over all DBs to allow client cross-compatibility
- Make using the API feel like you're using your own DB; endpoints shouldn't include usernames, simply names of contained buckets, completely separate from other users

## Installation

Assuming that you've correctly configured Go and have a MongoDB server listening on `:27017`, run

```shell script
$ go install github.com/fjah/bucketr
```

then simply execute the binary with

```shell script
$ MONGO_URI=mongodb://localhost:27017 bucketr
```

Note that Go and MongoDB need to be installed separately. In the future, this will all be done for you with Docker.