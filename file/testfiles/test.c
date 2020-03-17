#include <stdio.h>

int main(void) {
    int a = 0;
    char* hello_world = "Hello, how's it going?";
    printf("%s\n", hello_world);
    while (a == 0) {
        for (int i = 0; i < 100; i ++) {
            if (i % 2 == 0) {
                puts("FINE");
            } else {
                puts("NOT FINE");
            }
        }
    }
    return 0;
}