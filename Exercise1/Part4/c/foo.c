#include <pthread.h>
#include <stdio.h>

int i = 0;

// Note the return type: void*
void* incrementingThreadFunction(){
    for (int j = 1; i <= 100000; j++){
        i++;
    }
    return NULL;
}

void* decrementingThreadFunction(){
    for (int j = 1; i <= 100000; j++){
        i--;
    }
    return NULL;
}

//Build: gcc -pthread -o foo foo.c
int main(){
    // TODO: declare incrementingThread and decrementingThread (hint: google pthread_create)
    pthread_t incrementingThread;
    pthread_t decrementingThread;
    void *ret;
    pthread_create(&incrementingThread, NULL, incrementingThreadFunction, "Thread 1");
    pthread_create(&decrementingThread, NULL, decrementingThreadFunction, "Thread 2");
    
    pthread_join(incrementingThread, &ret);
    pthread_join(decrementingThread, &ret);
    
    printf("The magic number is: %d\n", i);
    return 0;
}
