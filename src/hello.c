#include <stdio.h>

int test() {
	printf("Hello World from Test\n");
	return 0;
}

int main(int argc, char argv)
{
	printf("Hello World\n");
	test();
	printf("Hello %s\n", "Daniel");
	return 0;
}