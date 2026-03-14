// MD5 基本関数
#[inline(always)]
fn f(x: u32, y: u32, z: u32) -> u32 { (x & y) | (!x & z) }
#[inline(always)]
fn g(x: u32, y: u32, z: u32) -> u32 { (x & z) | (y & !z) }
#[inline(always)]
fn h(x: u32, y: u32, z: u32) -> u32 { x ^ y ^ z }
#[inline(always)]
fn i(x: u32, y: u32, z: u32) -> u32 { y ^ (x | !z) }

// 左回転 (Left Rotate)
#[inline(always)]
fn lrot(x: u32, n: u32) -> u32 {
    (x << n) | (x >> (32 - n))
}

// ステップ実行用の補助関数
#[inline(always)]
fn step<F>(f_op: F, a: &mut u32, b: u32, c: u32, d: u32, x: u32, s: u32, t: u32) 
where F: Fn(u32, u32, u32) -> u32 {
    *a = a.wrapping_add(f_op(b, c, d)).wrapping_add(x).wrapping_add(t);
    *a = lrot(*a, s);
    *a = a.wrapping_add(b);
}

// 64個の定数 T
const T: [u32; 64] = [
    0xd76aa478, 0xe8c7b756, 0x242070db, 0xc1bdceee, 0xf57c0faf, 0x4787c62a, 0xa8304613, 0xfd469501,
    0x698098d8, 0x8b44f7af, 0xffff5bb1, 0x895cd7be, 0x6b901122, 0xfd987193, 0xa679438e, 0x49b40821,
    0xf61e2562, 0xc040b340, 0x265e5a51, 0xe9b6c7aa, 0xd62f105d, 0x02441453, 0xd8a1e681, 0xe7d3fbc8,
    0x21e1cde6, 0xc33707d6, 0xf4d50d87, 0x455a14ed, 0xa9e3e905, 0xfcefa3f8, 0x676f02d9, 0x8d2a4c8a,
    0xfffa3942, 0x8771f681, 0x6d9d6122, 0xfde5380c, 0xa4beea44, 0x4bdecfa9, 0xf6bb4b60, 0xbebfbc70,
    0x289b7ec6, 0xeaa127fa, 0xd4ef3085, 0x04881d05, 0xd9d4d039, 0xe6db99e5, 0x1fa27cf8, 0xc4ac5665,
    0xf4292244, 0x432aff97, 0xab9423a7, 0xfc93a039, 0x655b59c3, 0x8f0ccc92, 0xffeff47d, 0x85845dd1,
    0x6fa87e4f, 0xfe2ce6e0, 0xa3014314, 0x4e0811a1, 0xf7537e82, 0xbd3af235, 0x2ad7d2bb, 0xeb86d391,
];

pub fn md5_transform(state: &mut [u32; 4], block: &[u8]) {
    let mut a = state[0];
    let mut b = state[1];
    let mut c = state[2];
    let mut d = state[3];
    let mut x = [0u32; 16];

    // リトルエンディアンで4バイトずつ u32 に変換
    for i in 0..16 {
        x[i] = (block[i*4] as u32) |
               ((block[i*4+1] as u32) << 8) |
               ((block[i*4+2] as u32) << 16) |
               ((block[i*4+3] as u32) << 24);
    }

    // Round 1
    step(f, &mut a, b, c, d, x[0],  7, T[0]);  step(f, &mut d, a, b, c, x[1], 12, T[1]);
    step(f, &mut c, d, a, b, x[2], 17, T[2]);  step(f, &mut b, c, d, a, x[3], 22, T[3]);
    step(f, &mut a, b, c, d, x[4],  7, T[4]);  step(f, &mut d, a, b, c, x[5], 12, T[5]);
    step(f, &mut c, d, a, b, x[6], 17, T[6]);  step(f, &mut b, c, d, a, x[7], 22, T[7]);
    step(f, &mut a, b, c, d, x[8],  7, T[8]);  step(f, &mut d, a, b, c, x[9], 12, T[9]);
    step(f, &mut c, d, a, b, x[10], 17, T[10]); step(f, &mut b, c, d, a, x[11], 22, T[11]);
    step(f, &mut a, b, c, d, x[12],  7, T[12]); step(f, &mut d, a, b, c, x[13], 12, T[13]);
    step(f, &mut c, d, a, b, x[14], 17, T[14]); step(f, &mut b, c, d, a, x[15], 22, T[15]);

    // Round 2
    step(g, &mut a, b, c, d, x[1],  5, T[16]); step(g, &mut d, a, b, c, x[6],  9, T[17]);
    step(g, &mut c, d, a, b, x[11], 14, T[18]); step(g, &mut b, c, d, a, x[0], 20, T[19]);
    step(g, &mut a, b, c, d, x[5],  5, T[20]); step(g, &mut d, a, b, c, x[10],  9, T[21]);
    step(g, &mut c, d, a, b, x[15], 14, T[22]); step(g, &mut b, c, d, a, x[4], 20, T[23]);
    step(g, &mut a, b, c, d, x[9],  5, T[24]); step(g, &mut d, a, b, c, x[14],  9, T[25]);
    step(g, &mut c, d, a, b, x[3], 14, T[26]); step(g, &mut b, c, d, a, x[8], 20, T[27]);
    step(g, &mut a, b, c, d, x[13],  5, T[28]); step(g, &mut d, a, b, c, x[2],  9, T[29]);
    step(g, &mut c, d, a, b, x[7], 14, T[30]); step(g, &mut b, c, d, a, x[12], 20, T[31]);

    // Round 3
    step(h, &mut a, b, c, d, x[5],  4, T[32]); step(h, &mut d, a, b, c, x[8], 11, T[33]);
    step(h, &mut c, d, a, b, x[11], 16, T[34]); step(h, &mut b, c, d, a, x[14], 23, T[35]);
    step(h, &mut a, b, c, d, x[1],  4, T[36]); step(h, &mut d, a, b, c, x[4], 11, T[37]);
    step(h, &mut c, d, a, b, x[7], 16, T[38]); step(h, &mut b, c, d, a, x[10], 23, T[39]);
    step(h, &mut a, b, c, d, x[13],  4, T[40]); step(h, &mut d, a, b, c, x[0], 11, T[41]);
    step(h, &mut c, d, a, b, x[3], 16, T[42]); step(h, &mut b, c, d, a, x[6], 23, T[43]);
    step(h, &mut a, b, c, d, x[9],  4, T[44]); step(h, &mut d, a, b, c, x[12], 11, T[45]);
    step(h, &mut c, d, a, b, x[15], 16, T[46]); step(h, &mut b, c, d, a, x[2], 23, T[47]);

    // Round 4
    step(i, &mut a, b, c, d, x[0],  6, T[48]); step(i, &mut d, a, b, c, x[7], 10, T[49]);
    step(i, &mut c, d, a, b, x[14], 15, T[50]); step(i, &mut b, c, d, a, x[5], 21, T[51]);
    step(i, &mut a, b, c, d, x[12],  6, T[52]); step(i, &mut d, a, b, c, x[3], 10, T[53]);
    step(i, &mut c, d, a, b, x[10], 15, T[54]); step(i, &mut b, c, d, a, x[1], 21, T[55]);
    step(i, &mut a, b, c, d, x[8],  6, T[56]); step(i, &mut d, a, b, c, x[15], 10, T[57]);
    step(i, &mut c, d, a, b, x[6], 15, T[58]); step(i, &mut b, c, d, a, x[13], 21, T[59]);
    step(i, &mut a, b, c, d, x[4],  6, T[60]); step(i, &mut d, a, b, c, x[11], 10, T[61]);
    step(i, &mut c, d, a, b, x[2], 15, T[62]); step(i, &mut b, c, d, a, x[9], 21, T[63]);

    state[0] = state[0].wrapping_add(a);
    state[1] = state[1].wrapping_add(b);
    state[2] = state[2].wrapping_add(c);
    state[3] = state[3].wrapping_add(d);
}

pub fn md5_sum(msg: &[u8]) -> [u8; 16] {
    let mut state: [u32; 4] = [0x67452301, 0xefcdab89, 0x98badcfe, 0x10325476];
    let msg_len = msg.len();
    let mut i = 0;

    // 64バイトブロックの処理
    while i + 64 <= msg_len {
        md5_transform(&mut state, &msg[i..i+64]);
        i += 64;
    }

    // パディング処理
    let mut buffer = [0u8; 64];
    let remaining = msg_len - i;
    for j in 0..remaining { buffer[j] = msg[i + j]; }
    buffer[remaining] = 0x80;

    if remaining >= 56 {
        md5_transform(&mut state, &buffer);
        buffer = [0u8; 64];
    }

    // ビット長をリトルエンディアンで付与
    let bit_len = (msg_len as u64) * 8;
    for j in 0..8 {
        buffer[56 + j] = (bit_len >> (j * 8)) as u8;
    }
    md5_transform(&mut state, &buffer);

    let mut output = [0u8; 16];
    for j in 0..4 {
        output[j*4]   = (state[j] & 0xFF) as u8;
        output[j*4+1] = ((state[j] >> 8) & 0xFF) as u8;
        output[j*4+2] = ((state[j] >> 16) & 0xFF) as u8;
        output[j*4+3] = ((state[j] >> 24) & 0xFF) as u8;
    }
    output
}

fn main() {
    let msg = "abc";
    let hash = md5_sum(msg.as_bytes());
    print!("MD5(\"{}\") = ", msg);
    for b in hash {
        print!("{:02x}", b);
    }
    println!();
}

