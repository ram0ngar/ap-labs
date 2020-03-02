#include <stdio.h>

static char daytab[2][13] = {
  {0, 31, 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31},
  {0, 31, 29, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}
};

char *moth_name(int n) {
  static char *name[] = {
    "Illegal moth", "January", "February", "March", "April",
    "May", "June", "July", "August", "September", "October",
    "November", "December"
  };

  return (n < 1 || n > 12) ? name[0] : name[n];
}
void month_day(int year, int yearDay, int *pmonth, int *pday) {
  int i, leap;
  leap = year%4 == 0 && year%100 != 0 || year%400 == 0;
  for (i = 1; yearDay > daytab[leap][i]; i++){
    yearDay -= daytab[leap][i];
  }
    *pmonth = i;
    *pday = yearDay;
}

int main(int argc, char **argv) {
  // Place your magic here
  if (argc < 3||((int)atoi(argv[2])>366)) {
    printf("Not enough arguments o ilegal date\n");
    return -1;
  }
  int year,yearDay,month,day;
  year = (int)atoi(argv[1]);
  yearDay = (int)atoi(argv[2]);
  month_day(year, yearDay, &month, &day);
  char *pmonth = moth_name(month);
	if(day<10){
	  printf("%s 0%d, %d \n", pmonth, day, year);
	  return 0;
	}
	else{
	  printf("%s %d, %d \n", pmonth, day, year);
	  return 0;
	}
}
