#include <stdio.h>

struct list_head {
	struct list_head *next, *prev;
};

int main(void) {
    struct list_head m;
    puts("Hello world");
    return 0;
}