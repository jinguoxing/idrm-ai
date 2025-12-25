package main

import (
	"fmt"
	"idrm/pkg/validator"
)

// 示例结构体
type CreateUserRequest struct {
	Username string `json:"username" validate:"required,min=3,max=20,alphanum"`
	Password string `json:"password" validate:"required,min=6"`
	Email    string `json:"email" validate:"required,email"`
	Age      int    `json:"age" validate:"required,gte=18,lte=100"`
	Mobile   string `json:"mobile" validate:"required,mobile"`
	RealName string `json:"real_name" validate:"omitempty,chinese"`
}

func main() {
	// 初始化验证器（可选，会自动初始化）
	validator.Init()

	// 测试有效数据
	fmt.Println("=== 测试1: 有效数据 ===")
	validReq := CreateUserRequest{
		Username: "zhangsan",
		Password: "123456",
		Email:    "zhangsan@example.com",
		Age:      25,
		Mobile:   "13800138000",
		RealName: "张三",
	}

	if err := validator.Validate(validReq); err != nil {
		fmt.Printf("验证失败: %v\n", validator.GetErrorMsg(err))
	} else {
		fmt.Println("验证通过！")
	}

	// 测试无效数据
	fmt.Println("\n=== 测试2: 无效数据 ===")
	invalidReq := CreateUserRequest{
		Username: "ab",            // 太短
		Password: "123",           // 太短
		Email:    "invalid-email", // 格式错误
		Age:      15,              // 小于18
		Mobile:   "123",           // 格式错误
		RealName: "abc123",        // 包含非中文
	}

	if err := validator.Validate(invalidReq); err != nil {
		// 方法1: 获取详细错误字典
		fmt.Println("详细错误:")
		errMsgs := validator.GetErrorMsg(err)
		for field, msg := range errMsgs {
			fmt.Printf("  - %s: %s\n", field, msg)
		}

		// 方法2: 获取第一个错误
		fmt.Printf("\n第一个错误: %s\n", validator.GetFirstError(err))

		// 方法3: 格式化所有错误
		fmt.Printf("\n格式化错误: %s\n", validator.FormatError(err))

		// 方法4: 获取错误列表
		fmt.Println("\n错误列表:")
		for i, msg := range validator.GetErrorList(err) {
			fmt.Printf("  %d. %s\n", i+1, msg)
		}
	}

	// 测试自定义验证器
	fmt.Println("\n=== 测试3: 自定义验证器 ===")
	customReq := CreateUserRequest{
		Username: "testuser",
		Password: "password123",
		Email:    "test@example.com",
		Age:      20,
		Mobile:   "23456789012", // 不是1开头
		RealName: "测试用户",
	}

	if err := validator.Validate(customReq); err != nil {
		fmt.Printf("验证失败: %v\n", validator.GetErrorMsg(err))
	}
}
