package handler

import (
	"fmt"
	v1 "k8s.io/api/core/v1"
)

type CmdHandler struct {
}

func (c CmdHandler) OnAdd(obj interface{}, isInInitialList bool) {
	fmt.Println("add:", obj.(*v1.ConfigMap).Name)
}

func (c CmdHandler) OnUpdate(oldObj, newObj interface{}) {
	fmt.Println("update:", newObj.(*v1.ConfigMap).Name)
}

func (c CmdHandler) OnDelete(obj interface{}) {
	fmt.Println("delete:", obj.(*v1.ConfigMap).Name)
}
