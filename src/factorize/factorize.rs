use std::env;

/// 試し割り法による素因数分解
fn factorize(n: u64) {
    print!("{}: ", n);

    if n < 2 {
        println!("{}", n);
        return;
    }

    let mut temp = n;

    // 2で割れるだけ割る
    while temp % 2 == 0 {
        print!("2 ");
        temp /= 2;
    }

    // 3以上の奇数で試し割り（√temp まで）
    let mut i = 3;
    while i * i <= temp {
        while temp % i == 0 {
            print!("{} ", i);
            temp /= i;
        }
        i += 2;
    }

    // 最後に残った数が 2 以上なら素数
    if temp > 1 {
        print!("{}", temp);
    }
    println!();
}

fn main() {
    // 引数の取得
    let args: Vec<String> = env::args().collect();

    if args.len() < 2 {
        eprintln!("Usage: {} <num1> <num2> ... <numN>", args[0]);
        return;
    }

    // 第1引数はプログラム名なので飛ばす
    for arg in &args[1..] {
        // 文字列を u64 に変換 (エラーハンドリング)
        if let Ok(n) = arg.parse::<u64>() {
            if n > 0 {
                factorize(n);
            } else {
                println!("{}: Positive integer required.", arg);
            }
        } else {
            println!("{}: Invalid number.", arg);
        }
    }
}

