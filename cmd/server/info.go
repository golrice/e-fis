package main

const defaultBasePath = "/efis/"

type HttpInfo struct {
	addr     string
	basePath string
}

func NewHttpInfo(addr string) *HttpInfo {
	return &HttpInfo{
		addr:     addr,
		basePath: defaultBasePath,
	}
}
