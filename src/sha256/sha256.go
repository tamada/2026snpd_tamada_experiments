package main

import (
	"fmt"
)

// SHA-256 定数 K
var _K = [64]uint32{
	0x428a2f98, 0x71374491, 0xb5c0fbcf, 0xe9b5dba5, 0x3956c25b, 0x59f111f1, 0x923f82a4, 0xab1c5ed5,
	0xd807aa98, 0x12835b01, 0x243185be, 0x550c7dc3, 0x72be5d74, 0x80deb1fe, 0x9bdc06a7, 0xc19bf174,
	0xe49b69c1, 0xefbe4786, 0x0fc19dc6, 0x240ca1cc, 0x2de92c6f, 0x4a7484aa, 0x5cb0a9dc, 0x76f988da,
	0x983e5152, 0xa831c66d, 0xb00327c8, 0xbf597fc7, 0xc6e00bf3, 0xd5a79147, 0x06ca6351, 0x14292967,
	0x27b70a85, 0x2e1b2138, 0x4d2c6dfc, 0x53380d13, 0x650a7354, 0x766a0abb, 0x81c2c92e, 0x92722c85,
	0xa2bfe8a1, 0xa81a664b, 0xc24b8b70, 0xc76c51a3, 0xd192e819, 0xd6990624, 0xf40e3585, 0x106aa070,
	0x19a4c116, 0x1e376c08, 0x2748774c, 0x34b0bcb5, 0x391c0cb3, 0x4ed8aa4a, 0x5b9cca4f, 0x682e6ff3,
	0x748f82ee, 0x78a5636f, 0x84c87814, 0x8cc70208, 0x90befffa, 0xa4506ceb, 0xbef9a3f7, 0xc67178f2,
}

// 右回転 (Right Rotate)
func ror(x uint32, n uint) uint32 {
	return (x >> n) | (x << (32 - n))
}

// SHA-256 変換関数（1ブロック 64バイト処理）
func transform(state *[8]uint32, data []byte) {
	var w [64]uint32

	// メッセージスケジュール作成
	for i := 0; i < 16; i++ {
		w[i] = uint32(data[i*4])<<24 | uint32(data[i*4+1])<<16 |
			uint32(data[i*4+2])<<8 | uint32(data[i*4+3])
	}
	for i := 16; i < 64; i++ {
		g0 := ror(w[i-15], 7) ^ ror(w[i-15], 18) ^ (w[i-15] >> 3)
		g1 := ror(w[i-2], 17) ^ ror(w[i-2], 19) ^ (w[i-2] >> 10)
		w[i] = g1 + w[i-7] + g0 + w[i-16]
	}

	a, b, c, d, e, f, g, h := state[0], state[1], state[2], state[3], state[4], state[5], state[6], state[7]

	// メインループ
	for i := 0; i < 64; i++ {
		s1 := ror(e, 6) ^ ror(e, 11) ^ ror(e, 25)
		ch := (e & f) ^ (^e & g)
		t1 := h + s1 + ch + _K[i] + w[i]

		s0 := ror(a, 2) ^ ror(a, 13) ^ ror(a, 22)
		maj := (a & b) ^ (a & c) ^ (b & c)
		t2 := s0 + maj

		h, g, f, e, d, c, b, a = g, f, e, d+t1, c, b, a, t1+t2
	}

	state[0] += a
	state[1] += b
	state[2] += c
	state[3] += d
	state[4] += e
	state[5] += f
	state[6] += g
	state[7] += h
}

func Sum256(msg []byte) [32]byte {
	// 初期ハッシュ値
	state := [8]uint32{
		0x6a09e667, 0xbb67ae85, 0x3c6ef372, 0xa54ff53a,
		0x510e527f, 0x9b05688c, 0x1f83d9ab, 0x5be0cd19,
	}

	ln := len(msg)
	bitLen := uint64(ln) * 8

	// メインブロック処理
	i := 0
	for ; i+64 <= ln; i += 64 {
		transform(&state, msg[i:i+64])
	}

	// パディング処理
	tmp := make([]byte, 64)
	copy(tmp, msg[i:])
	tmp[ln-i] = 0x80

	if ln-i >= 56 {
		transform(&state, tmp)
		tmp = make([]byte, 64)
	}

	// 最後の8バイトに長さを書き込む
	for j := 0; j < 8; j++ {
		tmp[63-j] = byte(bitLen >> (j * 8))
	}
	transform(&state, tmp)

	// 出力の生成
	var out [32]byte
	for j := 0; j < 8; j++ {
		out[j*4] = byte(state[j] >> 24)
		out[j*4+1] = byte(state[j] >> 16)
		out[j*4+2] = byte(state[j] >> 8)
		out[j*4+3] = byte(state[j])
	}
	return out
}

func main() {
	input := "abc"
	hash := Sum256([]byte(input))

	fmt.Printf("Input:  %s\n", input)
	fmt.Printf("SHA-256: %x\n", hash)
}

