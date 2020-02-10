#include <stdio.h>

int add(int a, int b) {
    return a + b * a;
}

int main(void) {
    int x = 0; 
    for (int i = 0; i < 100; i += 1) {
        x = add(x, i);
        printf("%d\n", x);
    }
    return 0;
}