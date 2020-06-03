# Google Authenticator migration decoder

Convert [Google Authenticator](https://play.google.com/store/apps/details?id=com.google.android.apps.authenticator2)
[transfer links](https://github.com/google/google-authenticator-android/issues/118)
to plain [otpauth links](https://github.com/google/google-authenticator/wiki/Key-Uri-Format).

## Usage

Navigate to Menu -> Transfer accounts -> Export accounts
and extract the link from QR-code using your preferred software.

## Example

    go get github.com/dim13/otpauth
    otpauth -link "otpauth-migration://offline?data=CjEKCkhlbGxvId6tvu8SGEV4YW1wbGU6YWxpY2VAZ29vZ2xlLmNvbRoHRXhhbXBsZTAC"
    # Output:
    otpauth://totp/Example:alice@google.com?issuer=Example&secret=JBSWY3DPEHPK3PXP
