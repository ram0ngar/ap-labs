#define _XOPEN_SOURCE 500
#include <stdio.h>
#include <stdlib.h>
#include <errno.h>
#include <unistd.h>
#include <sys/types.h>
#include <linux/inotify.h>
#include <ftw.h>
#include <limits.h>
#include <string.h>
#include <stdint.h>
#include "logger.h"

#define EVENT_BUF_LEN     (10 * (sizeof(struct inotify_event) + NAME_MAX + 1))

int fd,wd;

static int inotifyAddPaths(const char *fpath, const struct stat *sb, int tflag, struct FTW *ftwbuf)
{
	/*adding the “/tmp” directory into watch list. Here, the suggestion is to validate the existence of the directory before adding into monitoring list.*/
    wd = inotify_add_watch( fd, fpath, IN_CREATE | IN_DELETE );
    if(wd==-1){
    	errorf("Error adding to iNotify\n");
    	exit(1);
    }
    return 0;           /* To tell nftw() to continue */
}


int main(int argc, char *argv[])
{
	char *i;
	ssize_t length;
  	char buffer[EVENT_BUF_LEN]__attribute__((aligned(8)));
  	char *path=calloc(30, sizeof(char));
  	int flags =0;
  	/*creating the INOTIFY instance*/
  	fd = inotify_init();

  	/*checking for error*/
  	if ( fd < 0 ) {
    	errorf( "inotify_init" );
    	exit(1);
  	}

  	flags = FTW_PHYS;

  	path=argv[1];

  	if (nftw((argc < 2) ? "." : path, inotifyAddPaths, 20, flags) == -1) {
        errorf("nftw");
        exit(1);
    }

  	printf("Monitoring directory: %s\n", path);


  	/*read to determine the event change happens on “/tmp” directory. Actually this read blocks until the change event occurs*/ 
  	for(;;){
	    length = read( fd, buffer, EVENT_BUF_LEN ); 
	  
	    /*checking for error*/
	    if ( length < 0 ) {
	      errorf( "read" );
	    }	  
	    if (length == 0) {
			panicf("read() from inotify fd returned 0!");
			exit(1);
		}
	    /*actually read return the list of change events happens. Here, read the change event one by one and process it accordingly.*/   
    	for(i=buffer;i<buffer+length;){
    		struct inotify_event *event = ( struct inotify_event * )i;     
	        if ( event->mask & IN_CREATE ) {
	          if ( event->mask & IN_ISDIR ) {
	            infof( "New directory %s created.\n", event->name );
	          }
	          else {
	            infof( "New file %s created.\n", event->name );
	          }
	        }
	        else if ( event->mask & IN_DELETE ) {
	          if ( event->mask & IN_ISDIR ) {
	            errorf( "Directory %s deleted.\n", event->name );
	          }
	          else {
	            errorf( "File %s deleted.\n", event->name );
	          }
	        }
	        printf("\n");

			if (event->len > 0)
				printf("        name = %s\n", event->name);
			i += sizeof(struct inotify_event) + event->len;
    	}
    }
    exit(0);
}
