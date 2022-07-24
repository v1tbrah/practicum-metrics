package api

import "strings"

type infoURL struct {
	typeM string
	nameM string
	valM  string
}

func newInfoUpdateURL(urlPath string) *infoURL {
	newInfoM := infoURL{}
	arrInfoM := strings.Split(strings.TrimPrefix(urlPath, "/update/"), "/")
	lenArrInfoM := len(arrInfoM)
	if lenArrInfoM > 0 {
		newInfoM.typeM = arrInfoM[0]
	}
	if lenArrInfoM > 1 {
		newInfoM.nameM = arrInfoM[1]
	}
	if lenArrInfoM > 2 {
		newInfoM.valM = arrInfoM[2]
	}
	return &newInfoM
}

func newInfoGetValueURL(urlPath string) *infoURL {
	newInfoM := infoURL{}
	arrInfoM := strings.Split(strings.TrimPrefix(urlPath, "/value/"), "/")
	lenArrInfoM := len(arrInfoM)
	if lenArrInfoM > 0 {
		newInfoM.typeM = arrInfoM[0]
	}
	if lenArrInfoM > 1 {
		newInfoM.nameM = arrInfoM[1]
	}
	return &newInfoM
}
