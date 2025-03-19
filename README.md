# Google Authenticator migration decoder

![Logo](images/otpauth.png)

Convert [Google Authenticator](https://play.google.com/store/apps/details?id=com.google.android.apps.authenticator2) `otpauth-migration://offline?data=...`
[transfer links](https://github.com/google/google-authenticator-android/issues/118)
to plain [otpauth links](https://github.com/google/google-authenticator/wiki/Key-Uri-Format).

## Usage

* Navigate to ⋮ → Transfer accounts → Export accounts.
* Extract migration link from QR-code using your preferred software.
* Pass link to `otpauth` tool.

### Flags

```
  -workdir string
      working directory to store eventual files (defaults to current one)
  -eval
      evaluate otps
  -http string
      serve http (e.g. :6060)
  -info
      display batch info
  -link string
      migration link (required)
  -qr
      generate QR-codes (optauth://)
  -rev
      reverse QR-code (otpauth-migration://)
  -file string
      input file with otpauth:// URLs (one per line)
  -migration-batch-img-prefix string
      prefix for batch QR code filenames (default "batch")
  -migration-batch-size int
      number of URLs to include in each batch (default: 7)
```

## Example

```
go install github.com/dim13/otpauth@latest
```

Or get latest binary [release](https://github.com/dim13/otpauth/releases/latest).

### Usage

```
~/go/bin/otpauth -link "otpauth-migration://offline?data=CjEKCkhlbGxvId6tvu8SGEV4YW1wbGU6YWxpY2VAZ29vZ2xlLmNvbRoHRXhhbXBsZTAC"
```

Will output:

```
otpauth://totp/Example:alice@google.com?issuer=Example&secret=JBSWY3DPEHPK3PXP
```

### QR-Codes

```
~/go/bin/otpauth -qr -link "otpauth-migration://offline?data=CjEKCkhlbGxvId6tvu8SGEV4YW1wbGU6YWxpY2VAZ29vZ2xlLmNvbRoHRXhhbXBsZTAC"
# view and scan *.png in current working directory
```

Will generate:

![Example](images/example.png)

### Process a File with otpauth URLs

You can also process a file containing multiple otpauth URLs (one per line) and generate QR codes for batches of 10 URLs:

```
~/go/bin/otpauth -file urls.txt -workdir output -migration-batch-img-prefix batch -migration-batch-size 10
```

This will:
1. Read all otpauth:// URLs from the file
2. Group them in batches of 10
3. Create migration payloads for each batch
4. Generate QR codes in the output directory with names like batch_1.png, batch_2.png, etc.

The generated QR codes can be scanned by Google Authenticator to import the accounts in each batch.

### Serve http
```
~/go/bin/otpauth -http=localhost:6060 -link "otpauth-migration://offline?data=CjEKCkhlbGxvId6tvu8SGEV4YW1wbGU6YWxpY2VAZ29vZ2xlLmNvbRoHRXhhbXBsZTAC"
```

Navigate to http://localhost:6060/

## Docker
A Docker container can also be used to run the application by building and running the image as following

#### Build image
From the current directory run
```
docker build . -t otpauth:latest
```

#### Run container
To start a container from the previously created image run
```
docker run --name otpauth -p 6060:6060 -v $(pwd)/workdir:/app/workdir --rm otpauth:latest -workdir /app/workdir -http :6060 -link "otpauth-migration://offline?data=CjEKCkhlbGxvId6tvu8SGEV4YW1wbGU6YWxpY2VAZ29vZ2xlLmNvbRoHRXhhbXBsZTAC"
```
```
-p 6060:6060
Map the host 6060 to the containr 6060

-v $(pwd)/workdir:/app/workdir
Map the host dir to the containr dir
```
Navigate to http://localhost:6060/

## Other projects

See also https://github.com/dim13/2fa for simple CLI 2FA evaluator
