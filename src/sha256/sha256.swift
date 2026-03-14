import Foundation

// SHA-256 定数 K
let K: [UInt32] = [
    0x428a2f98, 0x71374491, 0xb5c0fbcf, 0xe9b5dba5, 0x3956c25b, 0x59f111f1, 0x923f82a4, 0xab1c5ed5,
    0xd807aa98, 0x12835b01, 0x243185be, 0x550c7dc3, 0x72be5d74, 0x80deb1fe, 0x9bdc06a7, 0xc19bf174,
    0xe49b69c1, 0xefbe4786, 0x0fc19dc6, 0x240ca1cc, 0x2de92c6f, 0x4a7484aa, 0x5cb0a9dc, 0x76f988da,
    0x983e5152, 0xa831c66d, 0xb00327c8, 0xbf597fc7, 0xc6e00bf3, 0xd5a79147, 0x06ca6351, 0x14292967,
    0x27b70a85, 0x2e1b2138, 0x4d2c6dfc, 0x53380d13, 0x650a7354, 0x766a0abb, 0x81c2c92e, 0x92722c85,
    0xa2bfe8a1, 0xa81a664b, 0xc24b8b70, 0xc76c51a3, 0xd192e819, 0xd6990624, 0xf40e3585, 0x106aa070,
    0x19a4c116, 0x1e376c08, 0x2748774c, 0x34b0bcb5, 0x391c0cb3, 0x4ed8aa4a, 0x5b9cca4f, 0x682e6ff3,
    0x748f82ee, 0x78a5636f, 0x84c87814, 0x8cc70208, 0x90befffa, 0xa4506ceb, 0xbef9a3f7, 0xc67178f2
]

// 右回転 (Right Rotate)
func ror(_ x: UInt32, _ n: UInt32) -> UInt32 {
    return (x >> n) | (x << (32 - n))
}

func transform(state: inout [UInt32], data: ArraySlice<UInt8>) {
    var w = [UInt32](repeating: 0, count: 64)
    let dataArray = Array(data)

    // メッセージスケジュール作成
    for i in 0..<16 {
        w[i] = UInt32(dataArray[i*4]) << 24 | UInt32(dataArray[i*4+1]) << 16 |
               UInt32(dataArray[i*4+2]) << 8  | UInt32(dataArray[i*4+3])
    }
    for i in 16..<64 {
        let s0 = ror(w[i-15], 7) ^ ror(w[i-15], 18) ^ (w[i-15] >> 3)
        let s1 = ror(w[i-2], 17) ^ ror(w[i-2], 19) ^ (w[i-2] >> 10)
        w[i] = s1 &+ w[i-7] &+ s0 &+ w[i-16]
    }

    var (a, b, c, d, e, f, g, h) = (state[0], state[1], state[2], state[3], state[4], state[5], state[6], state[7])

    // メインループ
    for i in 0..<64 {
        let s1 = ror(e, 6) ^ ror(e, 11) ^ ror(e, 25)
        let ch = (e & f) ^ (~e & g)
        let t1 = h &+ s1 &+ ch &+ K[i] &+ w[i]
        
        let s0 = ror(a, 2) ^ ror(a, 13) ^ ror(a, 22)
        let maj = (a & b) ^ (a & c) ^ (b & c)
        let t2 = s0 &+ maj
        
        h = g; g = f; f = e; e = d &+ t1
        d = c; c = b; b = a; a = t1 &+ t2
    }

    state[0] &+= a; state[1] &+= b; state[2] &+= c; state[3] &+= d
    state[4] &+= e; state[5] &+= f; state[6] &+= g; state[7] &+= h
}

func sha256(msg: [UInt8]) -> [UInt8] {
    var state: [UInt32] = [
        0x6a09e667, 0xbb67ae85, 0x3c6ef372, 0xa54ff53a,
        0x510e527f, 0x9b05688c, 0x1f83d9ab, 0x5be0cd19
    ]

    let bitLen = UInt64(msg.count) * 8
    var i = 0
    while i + 64 <= msg.count {
        transform(state: &state, data: msg[i..<i+64])
        i += 64
    }

    // パディング処理
    var buffer = [UInt8](repeating: 0, count: 128)
    let remaining = msg.count - i
    for j in 0..<remaining { buffer[j] = msg[i + j] }
    buffer[remaining] = 0x80

    let padLen = (remaining >= 56) ? 128 : 64
    for j in 0..<8 {
        buffer[padLen - 1 - j] = UInt8((bitLen >> (j * 8)) & 0xFF)
    }

    transform(state: &state, data: buffer[0..<64])
    if padLen == 128 {
        transform(state: &state, data: buffer[64..<128])
    }

    // バイト配列に変換
    var out = [UInt8]()
    for s in state {
        out.append(UInt8((s >> 24) & 0xFF))
        out.append(UInt8((s >> 16) & 0xFF))
        out.append(UInt8((s >> 8) & 0xFF))
        out.append(UInt8(s & 0xFF))
    }
    return out
}

// 実行
let input = "abc"
let hash = sha256(msg: Array(input.utf8))
print("Input:  \(input)")
print("SHA-256: \(hash.map { String(format: "%02x", $0) }.joined())")

