package main

import (
	"fmt"
	"os"
	"os/exec"
)

//go:generate go build
func main() {
	// 加载命令行参数
	gen, err := Load()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	// 执行生成命令
	if err := gen.Generate(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	// 格式化
	cmd := exec.Command("go", "fmt", gen.Out)
	if err := cmd.Start(); err != nil {
		fmt.Printf("format go files failed,%v", err)
		os.Exit(1)
	}
	fmt.Printf(" ✅  完成任务\n\n")
}
