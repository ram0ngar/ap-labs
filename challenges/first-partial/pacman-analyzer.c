/*
Resources:
http://man7.org/linux/man-pages/man2/open.2.html
http://man7.org/linux/man-pages/man2/lseek.2.html
http://man7.org/linux/man-pages/man2/write.2.html
http://man7.org/linux/man-pages/man2/read.2.html

*/
#include <stdio.h>
#include <stdlib.h>
#include <sys/types.h>
#include <sys/stat.h>
#include <fcntl.h>
#include <unistd.h>
#include <string.h>
#include <stdbool.h>

#define INSTALLED 0
#define REMOVED 1
#define UPGRADED 2

#define REPORT_FILE "packages_report.txt"

struct Package{
    char removeDate[17];
    char installDate[17];
    char lastUpdate[17];
    char name[50];
    int updates;
    int status;
};
struct Hashtable{
    int size;
    int nelements;
    struct Package packeges[1000];
};

void analizeLog(char *logFile, char *report);
bool isAction(char c1, char c2);
void addToHashtable(struct Hashtable *ht, struct Package *p);
int getHashCode(char s[]);
bool findInHashtable(struct Hashtable *ht, char key[]);
struct Package *get(struct Hashtable *ht, char key[]);
void printHashtable(struct Hashtable *ht);
void htToString(char string[], struct Hashtable *ht);
void printPackage(struct Package *p);
void pToString(char string[], struct Package *ht);
void makeReport(char *reportS, int installed, int removed, int updated, int count, struct Hashtable *ht);


void addToHashtable(struct Hashtable *ht, struct Package *p){
    for (int i = 0; i < ht->nelements + 1; i++){
        int hashValue = getHashCode(p->name) + i;
        int index = hashValue % ht->size;
        if (strcmp(ht->packeges[index].name, "") == 0){
            ht->packeges[index] = *p;
            break;
        }
    }
    ht->nelements += 1;
}

bool isAction(char c1, char c2){
    if (c1 == 'i' || c1 == 'u'){
        return true;
    }
    else if (c1 == 'r' && c2 == 'e'){
        return true;
    }
    else{
        return false;
    }
}

int getHashCode(char s[]){
    int n = strlen(s);
    int hashValue = 0;

    for (int i = 0; i < n; i++){
        hashValue = hashValue * 31 + s[i];
    }

    hashValue = hashValue & 0x7fffffff;
    return hashValue;
}
bool findInHashtable(struct Hashtable *ht, char key[]){
    for (int i = 0; i < ht->nelements + 1; i++){
        int hashValue = getHashCode(key) + i;
        int index = hashValue % ht->size;
        if (strcmp(ht->packeges[index].name, key) == 0){
            return true;
        }
        else if (strcmp(ht->packeges[index].name, "") == 0){
            return false;
        }
    }
    return false;
}

struct Package *get(struct Hashtable *ht, char key[]){
    for (int i = 0; i < ht->nelements + 1; i++){
        int hashValue = getHashCode(key) + i;
        int index = hashValue % ht->size;
        if (strcmp(ht->packeges[index].name, key) == 0){
            return &ht->packeges[index];
        }
        else if (strcmp(ht->packeges[index].name, "") == 0){
            return NULL;
        }
    }
    return NULL;
}

void printHashtable(struct Hashtable *ht){
    printf("ht.size: %d\n", ht->size);
    printf("ht.nelements: %d\n", ht->nelements);
    printf("ht.array: \n");
    for (int i = 0; i < ht->size; i++){
        if (strcmp(ht->packeges[i].name, "") != 0){
            printPackage(&ht->packeges[i]);
            printf("\n");
        }
    }
}

void htToString(char string[], struct Hashtable *ht){
    for (int i = 0; i < ht->size; i++){
        if (strcmp(ht->packeges[i].name, "") != 0){
            pToString(string, &ht->packeges[i]);
            strcat(string, "\n\n");
        }
    }
}

void printPackage(struct Package *p){
    printf("  - Package name        : %s\n", p->name);
    printf("  - Install date        : %s\n", p->installDate);
    printf("  - Last update date    : %s\n", p->lastUpdate);
    printf("  - How many updates    : %d\n", p->updates);
    printf("  - Removal date        : %s\n", p->removeDate);
}

void pToString(char string[], struct Package *p){
    strcat(string, "  - Package name        : ");
    strcat(string, p->name);
    strcat(string, "\n");
    strcat(string, "  - Install date        : ");
    strcat(string, p->installDate);
    strcat(string, "\n");
    strcat(string, "  - Last update date    : ");
    strcat(string, p->lastUpdate);
    strcat(string, "\n");
    strcat(string, "  - How many updates    : ");
    char numBuf[20];
    sprintf(numBuf, "%d\n", p->updates);
    strcat(string, numBuf);
    strcat(string, "  - Removal date        : ");
    strcat(string, p->removeDate);
}

void makeReport(char *reportS, int installed, int removed, int updated, int count, struct Hashtable *ht){
    strcat(reportS, "Pacman Packages Report\n");
    strcat(reportS, "----------------------\n");
    char numBuf[120];
    sprintf(numBuf, "- Installed packages : %d\n- Removed packages   : %d\n- Upgraded packages  : %d\n- Current installed  : %d\n\n", installed, removed, updated, count);
    strcat(reportS, numBuf);
    htToString(reportS, ht);
}

int main(int argc, char **argv){

    if (argc < 2){
        printf("Usage:./pacman-analizer.o pacman.txt\n");
        return 1;
    }

    analizeLog(argv[1], REPORT_FILE);

    return 0;
}

void analizeLog(char *logFile, char *report){
    printf("Generating Report from: [%s] log file\n", logFile);

    struct Hashtable ht = {1000, 0};
    int installed = 0;
    int removed = 0;
    int updated = 0;
    int count = 0;

    int file = open(logFile, O_RDONLY);
    if (file == -1){
        printf("Unable to open the file\n");
        return;
    }
    int size = lseek(file, sizeof(char), SEEK_END);
   // close(file);
    file = open(logFile, O_RDONLY);
    if (file == -1){
        printf("Unable to open the file\n");
        return;
    }
    char buf[size];
    read(file, buf, size);
    close(file);
    buf[size - 1] = '\0';

    int indexDoc = 0;
    int state = 0;
    char date[17];
    char name[50];
    char action[10];
    bool validLine = false;
    while (indexDoc < size){
        if(state==0){ 
            if (buf[indexDoc] != 'f'){
                indexDoc++;
                int aux = 0;
                while (buf[indexDoc] != ']') {
                    date[aux] = buf[indexDoc];
                    aux++;
                    indexDoc++;
                }
                date[aux] = '\0';
                indexDoc = indexDoc + 2;
                state = 1
            ;
            }
            else
            {
                state = 1
            ;
            }
        
        }
        else if(state==1){
            while (buf[indexDoc] != ' '){
                indexDoc++;
            }
            indexDoc++;
            state = 2;
        
        }
        else if(state==2){
           int aux = 0;
            if (isAction(buf[indexDoc], buf[indexDoc + 1])){
                validLine = true;
                while (buf[indexDoc] != ' ') {
                    action[aux] = buf[indexDoc];
                    indexDoc++;
                    aux++;
                }
                action[aux] = '\0';
                indexDoc++;
                state = 3;
            }
            else
            {
                state = 4;
            }
            
        }
        else if(state==3){
           int aux = 0;
            while (buf[indexDoc] != ' '){
                name[aux] = buf[indexDoc];
                indexDoc++;
                aux++;
            }
            name[aux] = '\0';
            indexDoc++;
            state = 4;
        }
        else if(state==4){
            while (!(buf[indexDoc] == '\n' || buf[indexDoc] == '\0')){
                indexDoc++;
            }
            indexDoc++;
            if (validLine){
                if (!findInHashtable(&ht, name)){
                    struct Package p = {"", "", "", 0,'-', INSTALLED};
                    strcpy(p.name, name);
                    strcpy(p.installDate, date);
                    addToHashtable(&ht, &p);

                    installed++;
                }
                else
                {
                    struct Package *p1 = get(&ht, name);
                    if (strcmp(action, "installed") == 0){
                        if (p1->status == REMOVED){
                            p1->status = INSTALLED;
                            strcpy(p1->removeDate, "-");
                            removed--;
                        }
                    }
                    else if (strcmp(action, "removed") == 0){
                        if (p1->status == INSTALLED || p1->status == UPGRADED){
                            p1->status = REMOVED;
                            strcpy(p1->removeDate, date);
                            strcpy(p1->lastUpdate, date);
                            p1->updates = p1->updates + 1;
                            removed++;
                        }
                    }
                    else if (strcmp(action, "upgraded") == 0){
                        if (p1->status == INSTALLED){
                            p1->status = UPGRADED;
                            strcpy(p1->lastUpdate, date);
                            p1->updates = p1->updates + 1;
                            updated++;
                        }
                        else if (p1->status == UPGRADED){
                            strcpy(p1->lastUpdate, date);
                            p1->updates = p1->updates + 1;
                        }
                    }
                }
            }
            validLine = false;
            state = 0;
            if (indexDoc >= size - 1){
                indexDoc = indexDoc + 1;
            }
        }
    }
    count = installed - removed;
    char reportS[100000];
    makeReport(reportS, installed, removed, updated, count, &ht);
    file = open(report, O_CREAT | O_WRONLY, 0600);
    if (file == -1){
        printf("Failed open file\n");
        return;
    }
    write(file, reportS, strlen(reportS));
    close(file);

    printf("Report generated [%s]\n", report);
}
