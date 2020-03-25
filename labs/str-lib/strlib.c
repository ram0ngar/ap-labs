#include <stdio.h>
#include <stdlib.h>

int mystrlen(char *str){
	int count=0;
	char *p;
	p=str;
	while(*p!='\0'){
		p++;
		count++;
	}
    return count;
}
char *mystradd(char *origin, char *addition){
    int len=0;
    char *o1,*a1,*o2,*a2, *nWord;
    o1=o2=origin;
    a1=a2=addition;
    len=mystrlen(o1)+mystrlen(a1);
    char *output = malloc(sizeof(char)*len);
    nWord=output;
    while(*o2!='\0'){
    	*nWord=*o2;
    	o2++;
    	nWord++;    
    }
    while(*a2!='\0'){
    	*nWord=*a2;
    	a2++;
    	nWord++;
    }
    return output;
}
int mystrfind(char *origin, char *substr){
    int i;
    int aux = 0;
    int originLen = mystrlen(origin);
    int subLen = mystrlen(substr);
    for (i = 0; i <= originLen; i++) {
        if (origin[i] == substr[0] && aux < 1) {
            aux++;
        }
        else if (origin[i] == substr[aux]) {
            aux++;
        }
        else {
            aux = 0;
        }
        if (aux >= subLen) {
            return 1;
        }
    }
    return 0;
}