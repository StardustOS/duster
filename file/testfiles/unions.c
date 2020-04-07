#include <stdio.h>
union my_union
{
    int hello;
    char c;
};

int main(void) {
    union my_union m;
    m.hello = 5;
    puts("Hello world");
    return m.hello;
}