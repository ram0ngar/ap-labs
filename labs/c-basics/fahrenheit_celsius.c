#include <stdio.h>

#define   LOWER  atoi(argv[1])       /* lower limit of table */
#define   UPPER  atoi(argv[2])    //300     /* upper limit */
#define   STEP   atoi(argv[3])    //20      /* step size */

/* print Fahrenheit-Celsius table */

int main(int argc,char **argv)
{
    	int fahr=LOWER;

	//printf("counter: %3d",atoi(argv[1]));
	if(argc==4){
		for (fahr = LOWER; fahr <= UPPER; fahr = fahr + STEP)
		printf("Fahrenheit: %3d, Celcius: %6.1f\n", fahr, (5.0/9.0)*(fahr-32));
	}
	else{
		printf("F: %3d, C: %6.1f\n",fahr,(5.0/9.0)*(fahr-32));
	}
    return 0;
}
