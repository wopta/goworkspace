package lib

func GetContentType(ext string) string {
	m := make(map[string]string)
	m["doc"] = "application/msword"
	m["docx"] = "application/msword"
	m["pdf"] = "application/pdf"
	m["GIF"] = "image/gif"
	m["jpeg"] = "image/jpeg"
	m["jpg"] = "image/jpeg"
	m["jpe"] = "image/jpeg"
	m["PNG"] = "image/png"
	m["png"] = "image/png"
	m["tiff"] = "image/tiff"
	m["tif"] = "image/tiff"
	m["xls"] = "application/vnd.ms-excel"
	m["xlsx"] = "application/vnd.ms-excel"
	m["pptx"] = "application/vnd.ms-powerpoint"
	m["ppt"] = "application/vnd.ms-powerpoint"
	m["txt"] = "text/plain"
	m["zip"] = "application/zip"
	m["gzip"] = "application/x-gzip"
	return m[ext]
}
