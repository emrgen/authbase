package x

import v1 "github.com/emrgen/authbase/apis/v1"

type GetPage interface {
	GetPage() *v1.Page
}

func GetPageFromRequest(request GetPage) *v1.Page {
	page := v1.Page{
		Page: 0,
		Size: 20,
	}
	if request.GetPage() != nil {
		page.Page = request.GetPage().Page
		page.Size = request.GetPage().Size
	}

	return &page
}
