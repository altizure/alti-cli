package errors

const (
	// ErrNotImplemented is returned when the desired feature is not yet implemented.
	ErrNotImplemented AppError = "app: feature not implemented"
	// ErrNoConfig is returned when config file is not found.
	ErrNoConfig AppError = "app: no config"
	// ErrNotLogin is returned when user is not login.
	ErrNotLogin AppError = "app: not login"
	// ErrErrorCodeInvalid is returned when the input altizure error code is invalid.
	ErrErrorCodeInvalid AppError = "app: invalid error code"
	// ErrInvalidInput is returned when the input value is invalid.
	ErrInvalidInput AppError = "app: invalid input"
	// ErrProfileNotFound is returned when the queried profile is not found.
	ErrProfileNotFound ConfigError = "config: profile not found"
	// ErrProfileNotRemovable is returned when the default profile is chosen to be removed.
	ErrProfileNotRemovable ConfigError = "config: default profile not removable"
	// ErrClientInvisible is returned when the client is invisible to the api server.
	ErrClientInvisible ConfigError = "client: invisible"
	// ErrOffline is returned when the server is offline.
	ErrOffline ServerError = "server: offline"
	// ErrProjCreate is returned when a new project could not be created.
	ErrProjCreate ProjectError = "project: create"
	// ErrProjRemove is returned when a project could not be removed.
	ErrProjRemove ProjectError = "project: remove"
	// ErrProjNotFound is returned when a project is not found.
	ErrProjNotFound ProjectError = "project: project not found"
	// ErrImgNotFound is returned when an image could not be founded in the project.
	ErrImgNotFound ProjectError = "project: image not found"
	// ErrMetaNotFound is returned when a meta file could not be founded in the project.
	ErrMetaNotFound ProjectError = "project: meta file not found"
	// ErrMetaMisc is returned when a meta file is invalid or duplicated.
	ErrMetaMisc ProjectError = "project: meta file is invalid or duplicated"
	// ErrReportProj is returned when a project could not be reported.
	ErrReportProj ProjectError = "project: report error"
	// ErrFileNotImage is returned when a file is not a supported image.
	ErrFileNotImage FileError = "file: not image"
	// ErrFileNotZip is returned when a file is not a zip file.
	ErrFileNotZip FileError = "file: not zip"
	// ErrFileNotDir is returned when a file is not a directory.
	ErrFileNotDir FileError = "file: not directory"
	// ErrFileNotDirOrZip is returned when a file is not a directory and not a zip file.
	ErrFileNotDirOrZip FileError = "file: not directory or zip"
	// ErrFilesize is returned when the filesize of a file could not be determined.
	ErrFilesize FileError = "file: unknown filesize"
	// ErrFileImageDim is returned when the dimension of an image could not be determined.
	ErrFileImageDim FileError = "file: unknown image dimension"
	// ErrFileChecksum is returned when the checksum of a file could not be computed.
	ErrFileChecksum FileError = "file: unknown checksum"
	// ErrMetaFilenameInvalid is returned when the filename of meta file is invalid.
	ErrMetaFilenameInvalid FileError = "file: invalid meta filename"
	// ErrModelFilenameInvalid is returned when the filename of model file is invalid.
	ErrModelFilenameInvalid FileError = "file: invalid model filename"
	// ErrImgReg is returned when an image could not be registered for uploading.
	ErrImgReg UploadError = "upload: cannot register upload image"
	// ErrImgInvalid is returned when an image is regarded as invalid by the server.
	ErrImgInvalid UploadError = "upload: invalid image"
	// ErrClientTimeout is returned when the cli client could not get back Ready or Invalid image state within timeout.
	ErrClientTimeout UploadError = "upload: client timeout"
	// ErrUploadMethodInvalid is returned when the specified upload method is not supported.
	ErrUploadMethodInvalid UploadError = "upload: invalid upload method"
	// ErrNoBucketSuggestion is returned when no bucket suggestion is returned.
	ErrNoBucketSuggestion UploadError = "upload: no bucket suggestion"
	// ErrS3Error is returned when file upload operation could not result in ok status code.
	ErrS3Error UploadError = "upload: s3 error"
	// ErrMinioError is returned when file upload operation could not result in ok status code.
	ErrMinioError UploadError = "upload: minio error"
	// ErrBucketInvalid is returned when the provided bucket is invalid.
	ErrBucketInvalid UploadError = "upload: invalid bucket"
	// ErrNOSTS is returned when a new STS could not be obtained.
	ErrNOSTS UploadError = "upload: oss sts error"
	// ErrOSSUploaderNotFound is returned when an OSS uploader is not found.
	ErrOSSUploaderNotFound UploadError = "upload: oss uploader not found"
	// ErrImgMutateState is returned when the image state could not be mutated.
	ErrImgMutateState UploadError = "upload: cannot not mutate image state"
	// ErrModelMutateState is returned when the model state could not be mutated.
	ErrModelMutateState UploadError = "upload: cannot not mutate model state"
	// ErrModelReg is returned when a model could not be registered for uploading.
	ErrModelReg UploadError = "upload: cannot register upload model"
	// ErrMetaReg is returned when a meta file could not be registered for uploading.
	ErrMetaReg UploadError = "upload: cannot register meta file"
	// ErrMetaExisted is returned when a duplicated meta file is attempted to upload.
	ErrMetaExisted UploadError = "upload: meta file alreay existed"
	// ErrTaskStop is returned when a task could not be stopped
	ErrTaskStop TaskError = "task: task could not be stopped"
	// ErrTaskTypeInvalid is returned when the provided task type is invalid.
	ErrTaskTypeInvalid TaskError = "task: invalid task type"
	// ErrClientQuery is returned when the input gql query file is not found.
	ErrClientQuery ClientError = "client: query file not found"
	// ErrClientVar is returned when the input gql variable file is not found.
	ErrClientVar ClientError = "client: variable file not found"
	// ErrClientVarInvalid is returned when the input gql variable file is not valid.
	ErrClientVarInvalid ClientError = "client: variable file invalid"
	// ErrCurrencyInvalid is returned when the provided currency is invalid.
	ErrCurrencyInvalid BankError = "bank: invalid currency"
	// ErrTransferCoins is returned when the p2p coins give error.
	ErrTransferCoins BankError = "bank: transfer coins failed"
	// ErrInsufficientCoins is returned when the required coins is not enough.
	ErrInsufficientCoins BankError = "bank: insufficient coins"
)

// AppError is the application specific error.
type AppError string

func (e AppError) Error() string {
	return string(e)
}

// ConfigError is the config specific error.
type ConfigError string

func (e ConfigError) Error() string {
	return string(e)
}

// ServerError is the server specific error.
type ServerError string

func (e ServerError) Error() string {
	return string(e)
}

// ProjectError is the project related error.
type ProjectError string

func (e ProjectError) Error() string {
	return string(e)
}

// FileError is the file related error.
type FileError string

func (e FileError) Error() string {
	return string(e)
}

// UploadError is the upload related error.
type UploadError string

func (e UploadError) Error() string {
	return string(e)
}

// BankError is the ultitliy related error.
type BankError string

func (e BankError) Error() string {
	return string(e)
}

// TaskError is the task related error.
type TaskError string

func (e TaskError) Error() string {
	return string(e)
}

// ClientError is the client related error.
type ClientError string

func (e ClientError) Error() string {
	return string(e)
}
