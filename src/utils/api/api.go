package api

import "net/http"

func __internalGet(url string) (*http.Response,error) {
	resp,err := http.Get(url)
	if err!=nil{
		return nil,err
	}

	return resp,nil
}



func Get(url string) (*http.Response,error){
	return __internalGet(url)
}