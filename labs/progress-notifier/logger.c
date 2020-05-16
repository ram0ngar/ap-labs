#include <stdio.h>
#include <stdlib.h>
#include "logger.h"
#include <stdarg.h>
#include <syslog.h>
#include <signal.h>
#include <string.h>

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

int lType = 0;

void textcolor(int attr, int fg, int bg);
void textcolor(int attr, int fg, int bg)
{	char command[13];
	sprintf(command, "%c[%d;%d;%dm", 0x1B, attr, fg + 30, bg + 40);
	printf("%s", command);
}

void initLogger(char *logType) {
    printf("Initializing Logger on: %s\n", logType);
    if(strcmp("syslog",logType)==0 || strcmp("SYSLOG",logType)==0){
    	lType=1;
    }
    else{
    	lType=0;
    }
}

int infof(const char *format, ...) {
	int v;
	va_list ap;
	va_start(ap,format);
	if(lType==1){
		openlog("Advanced logger with syslog", LOG_CONS , LOG_LOCAL1);
		vsyslog(LOG_INFO,format,ap);
		closelog();
		va_end(ap);
		return v;
	}
	else{
		textcolor(BRIGHT, BLUE, BLACK);	
		printf("INFORMATION: ");
		v=vprintf(format,ap);
		textcolor(RESET, WHITE, BLACK);	
		va_end(ap);
		return v;
	}
    return 0;
}

int warnf(const char *format, ...) {
	int v;
	va_list ap;
	va_start(ap,format);
	if(lType==1){
		openlog("Advanced logger with syslog", LOG_CONS, LOG_LOCAL1);
		vsyslog(LOG_WARNING,format,ap);
		closelog();
		va_end(ap);
		return v;
	}
	else{
		textcolor(BRIGHT, YELLOW, BLACK);	
		printf("WARNING: ");
		v=vprintf(format,ap);
		textcolor(RESET, WHITE, BLACK);	
		va_end(ap);
		return v;
	}
    return 0;
}

int errorf(const char *format, ...) {
	int v;
	va_list ap;
	va_start(ap,format);
	if(lType==1){
		openlog("Advanced logger with syslog", LOG_CONS, LOG_LOCAL1);
		vsyslog(LOG_ERR,format,ap);
		closelog();
		va_end(ap);
		return v;
	}
	else{
		textcolor(BRIGHT, RED, BLACK);
		printf("ERROR: ");
		v=vprintf(format,ap);
		textcolor(RESET, WHITE, BLACK);	
		va_end(ap);
		return v;
	}
    return 0;
}
int panicf(const char *format, ...) {
	int v;
	va_list ap;
	va_start(ap,format);
	if(lType==1){
		openlog("Advanced logger with syslog", LOG_CONS, LOG_LOCAL1);
		vsyslog(LOG_CRIT,format,ap);
		closelog();
		va_end(ap);
		raise(SIGABRT);
		return v;
	}
	else{
		textcolor(BRIGHT, RED, BLACK);
		printf("PANIC: ");
		v=vprintf(format,ap);
		textcolor(RESET, WHITE, BLACK);	
		va_end(ap);
		raise(SIGABRT);
		return v;
	}
    return 0;
}