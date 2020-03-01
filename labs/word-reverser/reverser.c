#include <stdio.h>
#include <stdlib.h>
#define IN   1   /* inside a word */
#define OUT  0   /* outside a word */

void reverse(int c,int letters[]){
    int tmp;
    for(int i=0;i<c/2;i++){
        tmp=letters[i];
        letters[i]=letters[c-i-1];
        letters[c-i-1]=tmp;
    }
}
int main()

{
    int c, nl, nw, nc, state;
    state = OUT;
    int word[30];

    nl = nw = nc = 0;

    while ((c = getchar()) != EOF) {
	if (c == '\n'){
	    reverse(nc,word);
        for(int i=0;i<nc;i++){
            printf("%c",word[i]);
            word[i]=0;
        }
        
        nc=0;
        printf("\n");
        }
    else{
        //printf("%c",c);
        //printf("%i",nc);
        word[nc]=c;
        nc++;

    }
	/*if (c == ' ' || c == '\n' || c == '\t')
	    state = OUT;
	else if (state == OUT) {
	    state = IN;
	    ++nw;
	}*/
    }

    printf("Lines: %d,  Words: %d, Chars: %d\n", nl, nw, nc);
    return 0;
}
