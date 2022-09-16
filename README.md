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
```

## Example

**NOTE**: at least [Go](https://golang.org/dl/) 1.16 required, or use latest binary [release](https://github.com/dim13/otpauth/releases/latest).

```
go get github.com/dim13/otpauth
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
