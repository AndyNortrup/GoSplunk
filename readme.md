# GoSplunk
This is a Unofficial Splunk SDK for Go.  This SDK is intended to make it as simple as possible to create applications for Splunk using Go as possible.

Most of the existing work is focused on the ability to write a Modular Input.  For example to define a Splunk Modular Input scheme do the following:

```go
scheme := NewModInputScheme("Amazon S3",
  "Get data from Amazon S3.", true, StreamingModeXML)
scheme.AddArgument("name", "Resource name",
  "An S3 resource name without the leading s3://. "+
    "For example, for s3://bucket/file.txt specify bucket/file.txt. "+
    "You can also monitor a whole bucket (for example by specifying 'bucket'), "+
    "or files within a sub-directory of a bucket "+
    "(for example 'bucket/some/directory/'; note the trailing slash).",
  ModInputArgString, true, false)
scheme.AddArgument("key_id", "Key ID", "Your Amazon key ID.",
  ModInputArgString, true, false)
scheme.AddArgument("secret_key", "Secret key", "Your Amazon secret key.",
  ModInputArgString, true, false)

marshaledScheme, err := xml.Marshal(scheme)
```
