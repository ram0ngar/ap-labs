#include <stdio.h>
#include <stdlib.h>
#include <stdarg.h>
#include <signal.h>

#define RESET		0
#define BRIGHT 		1
#define DIM		2
#define UNDERLINE 	3
#define BLINK		4
#define REVERSE		7
#define HIDDEN		8

#define BLACK 		0
#define RED		1
#define GREEN		2
#define YELLOW		3
#define BLUE		4
#define MAGENTA		5
#define CYAN		6
#define	WHITE		7

void textcolor(int attr, int fg, int bg);
/*int main()
{	textcolor(REVERSE, RED, BLACK);	
	printf("In color\n");
	textcolor(RESET, WHITE, BLACK);	
	return 0;
}*/

void textcolor(int attr, int fg, int bg)
{	char command[13];

	/* Command is the control command to the terminal */
	sprintf(command, "%c[%d;%d;%dm", 0x1B, attr, fg + 30, bg + 40);
	printf("%s", command);
}

int infof(const char *format, ...){
	int v;
	va_list ap;
	va_start(ap,format);
	textcolor(BRIGHT, BLUE, BLACK);	
	printf("INFORMATION: ");
	v=vprintf(format,ap);
	textcolor(RESET, WHITE, BLACK);	
	va_end(ap);
	return v;
}

int warnf(const char *format, ...){
	int v;
	va_list ap;
	va_start(ap,format);
	textcolor(BRIGHT, YELLOW, BLACK);	
	printf("WARNING: ");
	v=vprintf(format,ap);
	textcolor(RESET, WHITE, BLACK);	
	va_end(ap);
	return v;
}

int errorf(const char *format, ...){
	int v;
	va_list ap;
	va_start(ap,format);
	textcolor(BRIGHT, RED, BLACK);
	printf("ERROR: ");
	v=vprintf(format,ap);
	textcolor(RESET, WHITE, BLACK);	
	va_end(ap);
	return v;
}

int panicf(const char *format, ...){
	int v;
	va_list ap;
	va_start(ap,format);
	textcolor(BRIGHT, RED, BLACK);
	printf("PANIC: ");
	v=vprintf(format,ap);
	textcolor(RESET, WHITE, BLACK);	
	va_end(ap);
	raise(SIGABRT);
	return v;
}
