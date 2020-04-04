#include<stdio.h>

struct k {
    int v;
    char c;
    float f;
};

struct m {
    size_t m;
    struct k meh;
    char b;
};

int main(void) {
    struct k my_struct;
    my_struct.c = 'c';
    my_struct.f = 0.1;
    struct m aNewStruct;
    return 0;
}