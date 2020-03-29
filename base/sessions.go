package base

type Page int

const (
	PageOne Page = iota + 1
	PageTwo
)

func GetCurrentPage(unqID string) (Page, error) {
	return PageOne, nil
}

func SetCurrentPage(unqID string, page Page) error {
	return nil
}
