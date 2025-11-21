package models

type CheckLinksRequest struct {
	Links []string `json:"links"`
}

type CheckLinksResponse struct {
	Links    map[string]string `json:"links"`
	LinksNum int               `json:"links_num"`
}

type LinkSet struct {
	ID    int
	Links map[string]string
}

type LinksListRequest struct {
	LinksList []int `json:"links_list"`
}
