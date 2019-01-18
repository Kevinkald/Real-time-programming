# Mutex and Channel basics

### What is an atomic operation?
> An atomic operation in concurrent programming are program operations that run completely independently
  of any other processes.
  An operation that can read and write on a single bus operation.

### What is a semaphore?
> A semaphore is a key giving the thread access to the resources. When the thread is done it needs to
  give away the key to the next thread. There may be many semaphores/keys at the same time.

### What is a mutex?
> Same as a semaphore, but now with only one key. The moderator decides who gets the key next.

### What is the difference between a mutex and a binary semaphore?
> A binary semaphore is used when there are only one shared resource. A mutex is a single key to the resource/resources.

### What is a critical section?
> A critical section is a section of resources that are protected from being accessed at the same time. Thus, threads are only able to use the resources one at a time.

### What is the difference between race conditions and data races?
 > A race condition occurs when two or more threads can access shared data and they try to change it at
   the same time. Because the thread scheduling algorithm can swap between threads at any time, you don't know the order in which the threads will attempt to access the shared data.

   A data race occurs when 2 instructions from different threads access the same memory location, at least one of these accesses is a write and there is no synchronization that is mandating any particular order among these accesses.

### List some advantages of using message passing over lock-based synchronization primitives.
> - No data races!
  - Messages are simpler than mutexes and semaphores.
  - More scalable

### List some advantages of using lock-based synchronization primitives over message passing.
> - More effective than message passing, better for miscrocontrollers etc.
  - Less overhead
  - Supported among the most programming languages.
