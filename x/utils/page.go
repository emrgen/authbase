package utils

import (
	v1 "github.com/emrgen/authbase/apis/v1"
)

type Pager interface {
	GetPage() *v1.Page
}

func GetPage(pager Pager) *v1.Page {
	page := pager.GetPage()
	if page == nil {
		page = &v1.Page{
			Page: 0,
			Size: 10,
		}
	}

	return page
}
