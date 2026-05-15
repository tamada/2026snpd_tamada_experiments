package main

import (
	"fmt"
	"os"
	"strconv"
)

// factorize は試し割り法で素因数分解を行う
func factorize(n int64) {
	fmt.Printf("%d: ", n)

	if n < 2 {
		fmt.Printf("%d\n", n)
		return
	}

	temp := n

	// 2で割れるだけ割る
	for temp%2 == 0 {
		fmt.Print("2 ")
		temp /= 2
	}

	// 3以上の奇数で試し割り（√temp まで）
	for i := int64(3); i*i <= temp; i += 2 {
		for temp%i == 0 {
			fmt.Printf("%d ", i)
			temp /= i
		}
	}

	// 最後に残った数が 2 以上なら素数
	if temp > 1 {
		fmt.Printf("%d", temp)
	}
	fmt.Println()
}

func main() {
	// 引数チェック
	if len(os.Args) < 2 {
		fmt.Printf("Usage: %s <num1> <num2> ... <numN>\n", os.Args[0])
		return
	}

	for _, arg := range os.Args[1:] {
		// 文字列を int64 に変換
		n, err := strconv.ParseInt(arg, 10, 64)
		if err != nil || n <= 0 {
			fmt.Printf("%s: Positive integer required.\n", arg)
			continue
		}
		factorize(n)
	}
}

