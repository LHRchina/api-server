package model

import "testing"

func TestGetAllUser(t *testing.T)  {
	ret, err := GetAllUser()
	t.Log(ret,err)
}
