Updates required:

* faith/set

  Changed Set new parameters

  Now:
    pushd ~/golang/src/gopkg.in/fatih/set.v0/
    git checkout 57907de300222151a123d29255ed17f5ed43fad3
  Future:
    Use new Set New interface

* https://github.com/twinj/uuid

  No longer has method Init
    - Is my simple comment out good enough?

  Now:

* Echo
    cannot use c.Response as type http.ResponseWriter



* Unix needs a manual go get
  go get golang.org/x/sys/unix



ulimit
  Change reference to this url  https://wilsonmar.github.io/maximum-limits/
  Covers all versions of mac

Max files
  We are concurrently opening all these files, they don't need to stay open.
  Open, read into memory, close - so we should not be hitting this limit.
