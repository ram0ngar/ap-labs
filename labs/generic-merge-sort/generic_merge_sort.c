#include <stdio.h>
#include <string.h>
#include <stdlib.h>


void mergeSort(void *array[], int leftPart, int rightPart, int (*comp)(void *, void *));
int numbercmp(char *chr, char *chr2);

int numbercmp(char *chr, char *chr2) {
  int a = atoi(chr);
  int b = atoi(chr2);
  if (a < b){
    return -1;
  }
  else if (a > b) {
    return 1;
  }
  else {
    return 0;
  }
}

void mergeSort(void *array[], int leftPart, int rightPart, int (*comp)(void *, void *)) {
  int i, j, k;

  if (leftPart < rightPart) {
    int middle = leftPart+(rightPart-leftPart)/2;
    mergeSort(array, leftPart, middle, comp);
    mergeSort(array, middle+1, rightPart, comp);
    int left = middle - leftPart + 1;
    int right = rightPart - middle;
    void *leftPartArr[left], *rightPartArr[right];
    for (i = 0; i < left; i++) {
      leftPartArr[i] = array[leftPart + i];
    }
    for (j = 0; j < right; j++) {
      rightPartArr[j] = array[middle + 1+ j];
    }
    i = 0;
    j = 0;
    k = leftPart;
    while (i < left && j < right) {
      if(comp(leftPartArr[i],rightPartArr[j]) < 1 || comp(leftPartArr[i], rightPartArr[j]) == 0) {
        array[k] = leftPartArr[i];
        i++;
      }
      else {
        array[k] = rightPartArr[j];
        j++;
      }
      k++;
    }
    while (i < left) {
      array[k] = leftPartArr[i];
      i++;
      k++;
    }
    while (j < right) {
      array[k] = rightPartArr[j];
      j++;
      k++;
    }
  }
}


int main(){
  char *num[] = {"78","98","0", "15","26"};
  char *str[] = {"hola", "pañal", "carne", "Tv", "Lupe", "Puerta"};
  int numLenght = sizeof(num)/sizeof(num[0]);
  int strLenght = sizeof(str)/sizeof(str[0]);
  mergeSort(num, 0, numLenght-1, numbercmp);
  for(int i = 0; i < numLenght; i++){
    printf("%s ", num[i]);
  }
  printf("\n");
  mergeSort(str, 0, strLenght-1, strcmp);
  for(int i = 0; i < strLenght; i++){
    printf("%s ", str[i]);
  }
  printf("\n");
  return 0;
}

