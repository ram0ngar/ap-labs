#include <stdio.h>
#include <stdlib.h>

int mystrlen(char *str);
char *mystradd(char *origin,char *addition);
int mystrfind(char *origin,char *substr);

int main(int argc,char **argv) {
	char *output;
	int len=mystrlen(argv[1]);
	printf("Initial length:\t%d\n",len);
	output = malloc(sizeof(char)*(len+mystrlen(argv[2])));
	output=mystradd(argv[1],argv[2]);
	printf("New string:\t%s\n",output);
	if(mystrfind(output,argv[3])){
		printf("Substring found:\tYes\n");
	}
	else{
		printf("Substring found:\tNo\n");
	}
    return 0;
}