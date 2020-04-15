#include <stdio.h>

int infof(const char *format, ...);
int warnf(const char *format, ...);
int errorf(const char *format, ...);
int panicf(const char *format, ...);

int main() {
	infof("Information test\n");
	warnf("Warning %s test\n","warnf");
	errorf("%s %s\n","Error","test");
	panicf("Panic test with core dump\n");
	return 0;
}
