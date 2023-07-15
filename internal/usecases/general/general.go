package general_uc

import (
	"github.com/nyelonong/boilerplate-go/internal/interfaces"
)

type Configuration struct {
	Upload       interfaces.Upload
	Notification interfaces.Notification
}

type GeneralUsecase struct {
	upload       interfaces.Upload
	notification interfaces.Notification
}

func New(config Configuration) interfaces.GeneralUC {
	return &GeneralUsecase{
		upload:       config.Upload,
		notification: config.Notification,
	}
}
