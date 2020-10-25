# Google Authenticator migration decoder

Convert [Google Authenticator](https://play.google.com/store/apps/details?id=com.google.android.apps.authenticator2)
[transfer links](https://github.com/google/google-authenticator-android/issues/118)
to plain [otpauth links](https://github.com/google/google-authenticator/wiki/Key-Uri-Format).

## Usage

* Navigate to ⋮ → Transfer accounts → Export accounts.
* Extract migration link from QR-code using your preferred software.
* Pass link to `otpauth` tool.

### Flags

```
  -eval
    	evaluate otps
  -link string
    	migration link (required)
  -qr
    	generate QR-codes
```

## Example

```
go get github.com/dim13/otpauth
otpauth -link "otpauth-migration://offline?data=CjEKCkhlbGxvId6tvu8SGEV4YW1wbGU6YWxpY2VAZ29vZ2xlLmNvbRoHRXhhbXBsZTAC"

# Output:
otpauth://totp/Example:alice@google.com?issuer=Example&secret=JBSWY3DPEHPK3PXP

# Output with QR code:
otpauth -link "otpauth-migration://offline?data=CjEKCkhlbGxvId6tvu8SGEV4YW1wbGU6YWxpY2VAZ29vZ2xlLmNvbRoHRXhhbXBsZTAC" -qr
# view and scan *.png in current working directory
```
