# Python 3.3.3 and 2.7.6
# python fo.py

from threading import Thread

# Potentially useful thing:
#   In Python you "import" a global variable, instead of "export"ing it when you declare it
#   (This is probably an effort to make you feel bad about typing the word "global")
i = 0
threading lock = Lock()

def incrementingFunction():
    global i
    for i in range 1000000:
        lock.acquire()
        i++
        lock.release()

def decrementingFunction():
    global i
    for i in range 1000000:
        lock.acquire()
        i--
        lock.release()

def main():
    global i

    incrementing = Thread(target = incrementingFunction, args = (),)
    decrementing = Thread(target = decrementingFunction, args = (),)
    
    # TODO: Start both threads
    
    incrementing.join()
    decrementing.join()
    
    print("The magic number is %d" % (i))


main()
