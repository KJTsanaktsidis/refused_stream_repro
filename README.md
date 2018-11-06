Reproduciton for REFUSED_STREAM errors in GCS library.

Build with `go build -o refused_stream .`
Running it for a few minutes, you should eventually start seeing errors like:

```
We got a close error: Post https://www.googleapis.com/upload/storage/v1/b/kjs_cool_vault_bucket/o?alt=json&prettyPrint=false&projection=full&uploadType=multipart: stream error: stream ID 45457; REFUSED_STREAM
```

because there are > MAX_CONCURRENT_STREAMS (100) open streams on the HTTP/2 connection.
