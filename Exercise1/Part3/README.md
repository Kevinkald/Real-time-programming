# Reasons for concurrency and parallelism


To complete this exercise you will have to use git. Create one or several commits that adds answers to the following questions and push it to your groups repository to complete the task.

When answering the questions, remember to use all the resources at your disposal. Asking the internet isn't a form of "cheating", it's a way of learning.

 ### What is concurrency? What is parallelism? What's the difference?
 > Concurrency refers to the ability to excecute programs/algorithms out of order and yet not affecting the final outcome.
   Paralellism is the ability to break programs/algorithms in to different parts and execute these parts simultaneously
   using two or more computing units, I.E cores. 
   Both concurrency and Parallellism can be used without the one or the other. Using parallellism is often preffered when
   solving tasks that are closely related to eachother, where the answers to each computation can be combined in to
   one general solution(Dynamic programming).
   Concurrent programming does not necessary address problems solving related tasks, however it supports
   pausing and continuing tasks.
 
 ### Why have machines become increasingly multicore in the past decade?
 > This is because of the reduction of transistors size and the ability to fit more transistors
   in to the circuit boards. The reduction of transistor size means that the clock speeds has increased
   quite dramatically the last decades but has slown down in the recent years. To continue increasing
   the number of computations one can achieve in a period of time one has to look in to
   multicore programming, where one utilized two or more cores for computation.
 
 ### What kinds of problems motivates the need for concurrent execution?
 (Or phrased differently: What problems do concurrency help in solving?)
 > Concurrent programming helps solving running different programs sharing the same resources.
   This technique involves using multiple threads. This technique makes the programmer
   believe that the programs execute "at the same time", yet they are not.
 
 ### Does creating concurrent programs make the programmer's life easier? Harder? Maybe both?
 (Come back to this after you have worked on part 4 of this exercise)
 > I would believe that concurrent programs makes the life easier for a programmer. It seems
   like it should provide an elegant solution to shared resources and how to run several
   programs at the same time using one one processor.
 
 ### What are the differences between processes, threads, green threads, and coroutines?
 > Processes: 		An instance of a computer program being executed. The process may include several threads
              		that execute instructions concurrently.
   Threads:   		Short for "thread of execution" is the way for a program to divide it self in to two
              		or more concurrent running tasks. Threads share the same resources.
   Green threads: 	Threads available to the programmer using a runtime library and not the 
                 	OS. They can therefore be used to simulate multithreading on systems
                 	that does not provide this function. Green threads = "User level threads". Virtual machine is given runtime from the OS and then the VM distributes the runtime to the threads.
   Coroutines:		Allows multiple entry points for a program supporting suspending and continuing execution at certain
   					locations in the running program. Thus this allows programs to jump from one location in the
   					program to perform another task then jump back and so forth. "Cooperation routines" lets the threads decide how the runtime is divided among threads.
 
 ### Which one of these do `pthread_create()` (C/POSIX), `threading.Thread()` (Python), `go` (Go) create?
 > These functions all create new threads of execution. 
   Go uses Goroutines
   C uses Regular threads
   Pyhton Green threads
 
 ### How does pythons Global Interpreter Lock (GIL) influence the way a python Thread behaves?
 > The GIL is a mutex protecting access to Python objects since Python is not thread safe and thus
   the mutex stops different threads from accessing the same resources at the same time, thus
   parallellism is not possible.
 
 ### With this in mind: What is the workaround for the GIL (Hint: it's another module)?
 > Using a multiprocess module.
 
 ### What does `func GOMAXPROCS(n int) int` change? 
 > It sets the maximum number of CPUs that can be running simultaneosuly and returns the previous 
   value.
