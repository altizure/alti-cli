package types

// ConvertToImageType converts file extension type to gql image type.
func ConvertToImageType(filetype string) string {
	ret := "JPEG"
	switch filetype {
	case "image/png":
		ret = "PNG"
	case "image/tiff":
		ret = "TIFF"
	case "image/jpge":
	default:
		ret = "JPEG"
	}
	return ret
}
