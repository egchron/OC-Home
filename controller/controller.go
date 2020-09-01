package controller

import (
	"okumoto.net/valvebox"
)

type Controller struct {
	ValveBox *valvebox.ValveBox
}

func (Controller) MainLoop() {
}
