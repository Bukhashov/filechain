package utils

var (
	imgAdmissible = []string {".png", ".jpg"}
	fileAdmissible = []string {".pdf"}
)

func ControlImgFormat(format string) bool {
	for i := range(imgAdmissible) {
		if format == imgAdmissible[i] {
			return true
		}
	}
	return false
}

func ControlFileFormat(format string) bool {
	for i := range(fileAdmissible) {
		if format == fileAdmissible[i] {
			return true
		}
	}
	return false
}

