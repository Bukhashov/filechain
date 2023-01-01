package folder

type Folder interface {
	New()
}
type folder struct {
	
}

func NewFolder() Folder {
	return &folder{

	}
}