#include <stdio.h>

struct k {
    int meh;
    int z;
};

int main(void) {
    int* pointer;
    struct k* m;
    printf("%p\n", pointer);
    printf("%p\n", m);
}