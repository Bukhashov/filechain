package utils

var (
	admissible = []string {".png", ".jpg"}
)

func ControlFormat(format string) bool {
	for i := range(admissible) {
		if format == admissible[i] {
			return true
		}
	}
	return false
}