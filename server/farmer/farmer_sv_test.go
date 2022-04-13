package farmer

import (
	"context"
	"encoding/json"
	"feimumoke/farming/v2/api/server"
	"fmt"
	"testing"
)

func Test_xxx(t *testing.T) {
	//Test()
	f := NewFarmerSv()
	ctx := context.Background()
	ctx = context.WithValue(ctx, "a", "ffvvv")
	val := `{"Name":"Zhangsan"}`
	for _, m := range server.FarmerService_ServiceDesc.Methods {
		if m.MethodName == "SelectFarmer" {
			resp, err := m.Handler(f, ctx, func(i interface{}) error {
				json.Unmarshal([]byte(val), i)
				return nil
			}, nil)
			fmt.Println(err)
			fmt.Println(resp)
		}
	}
}
