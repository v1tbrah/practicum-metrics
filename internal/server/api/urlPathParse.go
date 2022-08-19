package api

import (
	"strings"

	"github.com/rs/zerolog/log"
)

type pathInfo struct {
	typeM string
	nameM string
	valM  string
}

func newInfoUpdateURL(urlPath string) *pathInfo {
	log.Debug().Str("urlPath", urlPath).Msg("api.newInfoUpdateURL started")
	log.Debug().Msg("api.newInfoUpdateURL ended")

	newInfoM := pathInfo{}
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

func newInfoGetValueURL(urlPath string) *pathInfo {
	log.Debug().Str("urlPath", urlPath).Msg("api.newInfoGetValueURL started")
	log.Debug().Msg("api.newInfoGetValueURL ended")

	newInfoM := pathInfo{}
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
