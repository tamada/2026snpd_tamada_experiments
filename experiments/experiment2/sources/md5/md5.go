package main

import (
	"fmt"
)

// MD5 基本関数
func f(x, y, z uint32) uint32 { return (x & y) | (^x & z) }
func g(x, y, z uint32) uint32 { return (x & z) | (y & ^z) }
func h(x, y, z uint32) uint32 { return x ^ y ^ z }
func i(x, y, z uint32) uint32 { return y ^ (x | ^z) }

// 左回転 (Left Rotate)
func lrot(x uint32, n uint) uint32 {
	return (x << n) | (x >> (32 - n))
}

// ステップ実行用の補助関数
func step(fn func(uint32, uint32, uint32) uint32, a, b, c, d, x, s, t uint32) uint32 {
	a += fn(b, c, d) + x + t
	return lrot(a, uint(s)) + b
}

// 64個の定数 T
var t = [64]uint32{
	0xd76aa478, 0xe8c7b756, 0x242070db, 0xc1bdceee, 0xf57c0faf, 0x4787c62a, 0xa8304613, 0xfd469501,
	0x698098d8, 0x8b44f7af, 0xffff5bb1, 0x895cd7be, 0x6b901122, 0xfd987193, 0xa679438e, 0x49b40821,
	0xf61e2562, 0xc040b340, 0x265e5a51, 0xe9b6c7aa, 0xd62f105d, 0x02441453, 0xd8a1e681, 0xe7d3fbc8,
	0x21e1cde6, 0xc33707d6, 0xf4d50d87, 0x455a14ed, 0xa9e3e905, 0xfcefa3f8, 0x676f02d9, 0x8d2a4c8a,
	0xfffa3942, 0x8771f681, 0x6d9d6122, 0xfde5380c, 0xa4beea44, 0x4bdecfa9, 0xf6bb4b60, 0xbebfbc70,
	0x289b7ec6, 0xeaa127fa, 0xd4ef3085, 0x04881d05, 0xd9d4d039, 0xe6db99e5, 0x1fa27cf8, 0xc4ac5665,
	0xf4292244, 0x432aff97, 0xab9423a7, 0xfc93a039, 0x655b59c3, 0x8f0ccc92, 0xffeff47d, 0x85845dd1,
	0x6fa87e4f, 0xfe2ce6e0, 0xa3014314, 0x4e0811a1, 0xf7537e82, 0xbd3af235, 0x2ad7d2bb, 0xeb86d391,
}

func md5Transform(state *[4]uint32, block []byte) {
	a, b, c, d := state[0], state[1], state[2], state[3]
	var x [16]uint32

	// リトルエンディアンでの読み込み
	for j := 0; j < 16; j++ {
		x[j] = uint32(block[j*4]) | uint32(block[j*4+1])<<8 |
			uint32(block[j*4+2])<<16 | uint32(block[j*4+3])<<24
	}

	// Round 1
	a = step(f, a, b, c, d, x[0], 7, t[0]); d = step(f, d, a, b, c, x[1], 12, t[1])
	c = step(f, c, d, a, b, x[2], 17, t[2]); b = step(f, b, c, d, a, x[3], 22, t[3])
	a = step(f, a, b, c, d, x[4], 7, t[4]); d = step(f, d, a, b, c, x[5], 12, t[5])
	c = step(f, c, d, a, b, x[6], 17, t[6]); b = step(f, b, c, d, a, x[7], 22, t[7])
	a = step(f, a, b, c, d, x[8], 7, t[8]); d = step(f, d, a, b, c, x[9], 12, t[9])
	c = step(f, c, d, a, b, x[10], 17, t[10]); b = step(f, b, c, d, a, x[11], 22, t[11])
	a = step(f, a, b, c, d, x[12], 7, t[12]); d = step(f, d, a, b, c, x[13], 12, t[13])
	c = step(f, c, d, a, b, x[14], 17, t[14]); b = step(f, b, c, d, a, x[15], 22, t[15])

	// Round 2
	a = step(g, a, b, c, d, x[1], 5, t[16]); d = step(g, d, a, b, c, x[6], 9, t[17])
	c = step(g, c, d, a, b, x[11], 14, t[18]); b = step(g, b, c, d, a, x[0], 20, t[19])
	a = step(g, a, b, c, d, x[5], 5, t[20]); d = step(g, d, a, b, c, x[10], 9, t[21])
	c = step(g, c, d, a, b, x[15], 14, t[22]); b = step(g, b, c, d, a, x[4], 20, t[23])
	a = step(g, a, b, c, d, x[9], 5, t[24]); d = step(g, d, a, b, c, x[14], 9, t[25])
	c = step(g, c, d, a, b, x[3], 14, t[26]); b = step(g, b, c, d, a, x[8], 20, t[27])
	a = step(g, a, b, c, d, x[13], 5, t[28]); d = step(g, d, a, b, c, x[2], 9, t[29])
	c = step(g, c, d, a, b, x[7], 14, t[30]); b = step(g, b, c, d, a, x[12], 20, t[31])

	// Round 3
	a = step(h, a, b, c, d, x[5], 4, t[32]); d = step(h, d, a, b, c, x[8], 11, t[33])
	c = step(h, c, d, a, b, x[11], 16, t[34]); b = step(h, b, c, d, a, x[14], 23, t[35])
	a = step(h, a, b, c, d, x[1], 4, t[36]); d = step(h, d, a, b, c, x[4], 11, t[37])
	c = step(h, c, d, a, b, x[7], 16, t[38]); b = step(h, b, c, d, a, x[10], 23, t[39])
	a = step(h, a, b, c, d, x[13], 4, t[40]); d = step(h, d, a, b, c, x[0], 11, t[41])
	c = step(h, c, d, a, b, x[3], 16, t[42]); b = step(h, b, c, d, a, x[6], 23, t[43])
	a = step(h, a, b, c, d, x[9], 4, t[44]); d = step(h, d, a, b, c, x[12], 11, t[45])
	c = step(h, c, d, a, b, x[15], 16, t[46]); b = step(h, b, c, d, a, x[2], 23, t[47])

	// Round 4
	a = step(i, a, b, c, d, x[0], 6, t[48]); d = step(i, d, a, b, c, x[7], 10, t[49])
	c = step(i, c, d, a, b, x[14], 15, t[50]); b = step(i, b, c, d, a, x[5], 21, t[51])
	a = step(i, a, b, c, d, x[12], 6, t[52]); d = step(i, d, a, b, c, x[3], 10, t[53])
	c = step(i, c, d, a, b, x[10], 15, t[54]); b = step(i, b, c, d, a, x[1], 21, t[55])
	a = step(i, a, b, c, d, x[8], 6, t[56]); d = step(i, d, a, b, c, x[15], 10, t[57])
	c = step(i, c, d, a, b, x[6], 15, t[58]); b = step(i, b, c, d, a, x[13], 21, t[59])
	a = step(i, a, b, c, d, x[4], 6, t[60]); d = step(i, d, a, b, c, x[11], 10, t[61])
	c = step(i, c, d, a, b, x[2], 15, t[62]); b = step(i, b, c, d, a, x[9], 21, t[63])

	state[0] += a; state[1] += b; state[2] += c; state[3] += d
}

func md5Sum(msg []byte) [16]byte {
	state := [4]uint32{0x67452301, 0xefcdab89, 0x98badcfe, 0x10325476}
	msgLen := len(msg)
	
	// 64バイトごとのメイン処理
	i := 0
	for ; i+64 <= msgLen; i += 64 {
		md5Transform(&state, msg[i:i+64])
	}

	// パディング処理
	buffer := make([]byte, 64)
	remaining := msgLen - i
	copy(buffer, msg[i:])
	buffer[remaining] = 0x80

	if remaining >= 56 {
		md5Transform(&state, buffer)
		for j := range buffer { buffer[j] = 0 }
	}

	// 長さをリトルエンディアンでセット (bit単位)
	bitLen := uint64(msgLen) * 8
	for j := 0; j < 8; j++ {
		buffer[56+j] = byte(bitLen >> (j * 8))
	}
	md5Transform(&state, buffer)

	// 結果をバイト配列に書き出し
	var output [16]byte
	for j := 0; j < 4; j++ {
		output[j*4+0] = byte(state[j])
		output[j*4+1] = byte(state[j] >> 8)
		output[j*4+2] = byte(state[j] >> 16)
		output[j*4+3] = byte(state[j] >> 24)
	}
	return output
}

func main() {
	msg := "abc"
	result := md5Sum([]byte(msg))
	fmt.Printf("MD5(\"%s\") = %x\n", msg, result)
}

