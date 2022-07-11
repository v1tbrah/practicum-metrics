package metric

import "strings"

type Info struct {
	typeM string
	nameM string
	valM  string
}

func (i *Info) TypeM() string {
	return i.typeM
}

func (i *Info) NameM() string {
	return i.nameM
}

func (i *Info) ValM() string {
	return i.valM
}

func NewInfoFromUpdateURL(urlPath string) *Info {
	newInfoM := Info{}
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

func NewInfoFromGetValueURL(urlPath string) *Info {
	newInfoM := Info{}
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
