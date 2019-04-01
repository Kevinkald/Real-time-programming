# Elevator project

## How to run

The program can be run by typing "go run main.go" followed by the ID and elevatorServer/simulator port for the elevator. For an elevator with ID = 2 and an ElevatorServer listening for commands at localhost:port 15657, write:

  `go run main.go 2 15657`

If you are using the binary executable file and would like to start three elevator instances using the simulator, write the following in three separate terminal windows:
  
  `./main 1 16001`
  `./main 2 16002`
  `./main 3 16003`
  
## Our solution

For this project we chose the programming language Go, which utilizes channels for communcation between modules. The concept was fairly easy to understand and also a very powerful tool for this task.

The solution builds on a peer to peer logic, where all nodes/elevators alive on the network shares all information about each other. The data structure containing all this information is a map keyed on the unique elevator ID's (1,2,3..). The value of the corresponding key is a struct of the elevator state information and order matrix. 

For the elevators to have the same information about each other, broadcasting to port using UDP was utilized. When an elevator receives a message from another elevator, it updates the corresponding elevator state information in its map. The orders are then synchronized by iterating through both maps and taking union of the elements. This makes sure that all new orders are updated in all maps. To remove orders the synchronization algorithm checks if two elements are different, and if that is the case it checks if theres an elevator standing in that corresponding floor with the doors open. If that is the case, we can safely remove that order since it has been served.

When a button is pushed, the order is directly fed in to a delegation algorithm which has access to all elevator information using the map described above. The algorithm only delegates to alive elevators. When an elevator is chosen, the corresponding queue element in the map keyed on the chosen elevator is updated accordingly. This is then broadcasted and thus the change is picked up by the chosen elevator. Cab calls are obviously directly given to the elevator polling the buttons.

When an elevator loses network connection, all elevators notice that theres an change in the number of elevators alive on the network. All hall orders are then redistributed among the elevators that the corresponding elevator sees as alive. Thus making sure that no orders are lost. When an elevator enters the network, the same procedure is applied as well as the elevator gets its cab calls back.

If an elevator loses motor power, a timeout is given if the elvator is in a moving state but does not change floor for 5 seconds. After the timeout the code restarts and starts an initializing procedure, which will not continue until power is restored.

## Modules

### Finite State Machine Module

The finite state machine (fsm) module implements single elevator control through different states and events. The module receives local orders from the queuedistribution module and executes them in compliance to the requirements specification. This includes handling order execution, the doorlight, motor direction and engine errors.

### Queue Distributor Module

The queue distributor module has a multitude of responsibilites concerning the synchronization of elevators. This includes storing information about the position, direction, state and orders of all the elevators in the system. The module is also responsible for distributing received orders to elevators in the system based on a costfunction.

### Network Module

The network module being used in this project was handed out and can be found [here](https://github.com/TTK4145/Network-go). The code is mostly untouched, with only small changes to how the list of alive nodes is handled when the system loses connection to the router. The features of this module are as described below:

Channel-in/channel-out pairs of (almost) any custom or built-in datatype can be supplied to a pair of transmitter/receiver functions. Data sent to the transmitter function is automatically serialized and broadcasted on the specified port. Any messages received on the receiver's port are deserialized (as long as they match any of the receiver's supplied channel datatypes) and sent on the corresponding channel. See bcast.Transmitter and bcast.Receiver.

Peers on the local network can be detected by supplying your own ID to a transmitter and receiving peer updates (new, current and lost peers) from the receiver. See peers.Transmitter and peers.Receiver.

### Elevator I/O

The Elevator I/O module for this project was handed out and can be found [here](https://github.com/TTK4145/driver-go). The module acts as an interface between the elevator hardware and and our code, and is responsible for setting the motor direction, updating lights and polling buttons and floor sensors.

### Simulator v2

The Simulator module for this project was handed out and can be found [here](https://github.com/TTK4145/Simulator-v2). The module acts as a replacement for the physical elevator hardware.

### Elevator server

The elevator server for this project was handed out and can be found [here](https://github.com/TTK4145/elevator-server). This module connects the physical elevator hardware with the software using TCP. This server can be swapped with the simulator mentioned above.
