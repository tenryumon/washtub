package general

import (
	"github.com/nyelonong/boilerplate-go/internal/interfaces"
)

type Configuration struct {
	AuthorizationUC interfaces.AuthorizationUC
	GeneralUC       interfaces.GeneralUC
}

type GeneralView struct {
	authorizationUC interfaces.AuthorizationUC
	generalUC       interfaces.GeneralUC
}

func New(config Configuration) *GeneralView {
	return &GeneralView{
		authorizationUC: config.AuthorizationUC,
		generalUC:       config.GeneralUC,
	}
}
