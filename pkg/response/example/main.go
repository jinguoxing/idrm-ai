package main

import (
	"fmt"
	"idrm/pkg/errorx"
	"idrm/pkg/response"
	"idrm/pkg/validator"
	"net/http/httptest"
)

type CreateUserRequest struct {
	Name  string `json:"name" validate:"required,min=2,max=50"`
	Email string `json:"email" validate:"required,email"`
	Age   int    `json:"age" validate:"required,gte=18,lte=100"`
}

func main() {
	fmt.Println("=== Response 包使用示例 ===\n")

	// 示例1: 成功响应
	fmt.Println("1. 成功响应:")
	w1 := httptest.NewRecorder()
	data := map[string]interface{}{
		"id":   1,
		"name": "测试用户",
	}
	response.Success(w1, data)
	fmt.Printf("Status: %d\n", w1.Code)
	fmt.Printf("Body: %s\n\n", w1.Body.String())

	// 示例2: 验证错误响应
	fmt.Println("2. 验证错误响应:")
	w2 := httptest.NewRecorder()
	req := CreateUserRequest{
		Name:  "a",
		Email: "invalid-email",
		Age:   15,
	}
	if err := validator.Validate(req); err != nil {
		errMsgs := validator.GetErrorMsg(err)
		response.ErrorValidation(w2, errMsgs)
	}
	fmt.Printf("Status: %d\n", w2.Code)
	fmt.Printf("Body: %s\n\n", w2.Body.String())

	// 示例3: 自定义详细错误
	fmt.Println("3. 详细错误响应:")
	w3 := httptest.NewRecorder()
	response.ErrorDetailed(
		w3,
		"idrm.category.duplicate_code",
		"类别代码已存在",
		"请使用不同的类别代码",
		"数据库唯一索引冲突",
		map[string]string{
			"duplicate_code": "TEST001",
			"field":          "code",
		},
	)
	fmt.Printf("Status: %d\n", w3.Code)
	fmt.Printf("Body: %s\n\n", w3.Body.String())

	// 示例4: 404错误
	fmt.Println("4. NotFound错误:")
	w4 := httptest.NewRecorder()
	response.NotFound(w4, "类别")
	fmt.Printf("Status: %d\n", w4.Code)
	fmt.Printf("Body: %s\n\n", w4.Body.String())

	// 示例5: 分页响应
	fmt.Println("5. 分页响应:")
	w5 := httptest.NewRecorder()
	list := []map[string]interface{}{
		{"id": 1, "name": "项目1"},
		{"id": 2, "name": "项目2"},
	}
	response.SuccessPage(w5, list, 100, 1, 10)
	fmt.Printf("Status: %d\n", w5.Code)
	fmt.Printf("Body: %s\n\n", w5.Body.String())

	// 示例6: 与errorx集成
	fmt.Println("6. ErrorX集成:")
	w6 := httptest.NewRecorder()
	err := errorx.NewWithMsg(400, "参数错误")
	response.Error(w6, err)
	fmt.Printf("Status: %d\n", w6.Code)
	fmt.Printf("Body: %s\n\n", w6.Body.String())
}
