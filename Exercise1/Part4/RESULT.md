Golang:
	The "magic number" seems to be rather random, taking values between [-100 000, 100 000]. This is probably
	because the two threads are sharing the resource of the global variable i, and thus the two threads
	are competing against each other and therefore sometimes fails to increment/decrement i.

C:
	Running the code using: gcc -pthread -o foo foo.c ./foo
	gave the result: 	The magic number is: 2147483647
	It looks like the first thread stole the access to the resource variable i and then
	maxed out the integer i.

Python:
	Running the code using: chmod +x foo.py ./foo.py
	gave the result: 
