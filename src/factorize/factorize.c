#include <stdio.h>
#include <stdlib.h>

/**
 * 試し割り法による素因数分解
 * @param n 分解する数値
 */
void factorize(long long n) {
    printf("%lld: ", n);
    
    if (n < 2) {
        printf("%lld\n", n);
        return;
    }

    // 2で割れるだけ割る
    while (n % 2 == 0) {
        printf("2 ");
        n /= 2;
    }

    // 3以上の奇数で試し割り（√n まで）
    for (long long i = 3; i * i <= n; i += 2) {
        while (n % i == 0) {
            printf("%lld ", i);
            n /= i;
        }
    }

    // 最後に残った数が 2 以上なら、それは素数
    if (n > 1) {
        printf("%lld", n);
    }
    printf("\n");
}

int main(int argc, char *argv[]) {
    if (argc < 2) {
        printf("Usage: %s <num1> <num2> ... <numN>\n", argv[0]);
        return 1;
    }

    for (int i = 1; i < argc; i++) {
        // 文字列を数値に変換 (long long型)
        long long n = atoll(argv[i]);
        if (n <= 0) {
            printf("%s: Positive integer required.\n", argv[i]);
            continue;
        }
        factorize(n);
    }

    return 0;
}

