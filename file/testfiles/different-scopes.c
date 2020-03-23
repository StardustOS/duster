#include <stdio.h>

int main(void) {
    int k = 0;
    for (int i = 0; i < 100; i += 1) {
        int factor = 0;
        k += 1;
        factor = 100 * k;
        if (k > 100) {
            int j = 100;
            k += j;
        }
        k += 2;
    }
    return 0;
}