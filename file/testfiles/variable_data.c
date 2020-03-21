#include <stdio.h>

// int y;

// struct k {
//     int v;
//     char c;
//     float f;
// };

// int clean(int val1, int val2) {
//     int a = val1 + val2 + y;
//     struct k l;
//     l.v = 1;
//     printf("%d\n", a);
//     int y = 20;
//     printf("%d\n", y);
// }

int main(void) {
    int a = 0;
    char* hello_world = "Hello, how's it going?";
    printf("%s\n", hello_world);
    while (a == 0) {
        for (int i = 0; i < 100; i ++) {
            if (i % 2 == 0) {
                puts("FINE");
                a = 1;
            } else {
                puts("NOT FINE");
            }
        }
    }
    int meh = 100;
    return 0;
}