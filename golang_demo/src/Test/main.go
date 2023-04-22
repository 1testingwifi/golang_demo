package main

import "fmt"

type FixedData [10]byte

func main() {
	var data FixedData
	str := "你好，世界！"
	copy(data[:], str[:9]) //定义长度
	fmt.Println(string(data[:]))
}
