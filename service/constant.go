package service

// NormalMode is the literal of normal mode returned from gql.
const NormalMode = "Normal"

// ReadOnlyMode is the literal of read-only mode returned from gql.
const ReadOnlyMode = "ReadOnly"

// DirectUploadMethod is the literal used in the arags of the import command.
const DirectUploadMethod = "direct"

// S3UploadMethod is the literal used in the arags of the import command.
const S3UploadMethod = "s3"

// MinioUploadMethod is the literal used in the arags of the import command.
const MinioUploadMethod = "minio"

// OSSUploadMethod is the literal used in the arags of the import command.
const OSSUploadMethod = "oss"

// Pending represents the image or model or meta pending state.
const Pending = "Pending"

// Ready represents the image or model or meta ready state.
const Ready = "Ready"

// Failed represents the image or model or meta failed state.
const Failed = "Failed"

// Yes represents an answer of yes.
const Yes = "YES"

// ValidMetafileNames specifies the valid metafile names.
var ValidMetafileNames = []string{"camera.txt", "pose.txt", "group.txt", "initial.xms", "initial.xms.zip"}
