package v2

import (
	"github.com/KaySar12/NextZenOS/codegen"
	"github.com/KaySar12/NextZenOS/service"
)

type CasaOS struct {
	fileUploadService *service.FileUploadService
}

func NewCasaOS() codegen.ServerInterface {
	return &CasaOS{
		fileUploadService: service.NewFileUploadService(),
	}
}
