package main

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	password := "123456"
	
	// 生成哈希
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Printf("Error generating hash: %v\n", err)
		return
	}
	fmt.Printf("Generated hash: %s\n", string(hashedPassword))
	fmt.Printf("Hash length: %d\n", len(hashedPassword))
	
	// 验证密码
	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		fmt.Printf("Password verification failed: %v\n", err)
	} else {
		fmt.Println("Password verification SUCCESS!")
	}
	
	// 验证错误密码
	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte("wrongpassword"))
	if err != nil {
		fmt.Printf("Wrong password verification correctly failed: %v\n", err)
	}
}
