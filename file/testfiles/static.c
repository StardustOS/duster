#include <stdio.h>

int my_func(int c) {
	static int k = 0;
	return c + k;
}

int main(void) {
    //static int my_integer;
    puts("Hello world");
    int c = my_func(150);
    //my_integer = 100;
    return 0;
}
