package fsmdummy

import "../elevio"
import "fmt"
import "../../variabletypes"
import "../../config"


func FsmDummy(){

    elevio.Init(config.SimulatorPort)
    
    var d variabletypes.MotorDirection = variabletypes.MD_Up
    elevio.SetMotorDirection(d)
    
    //drv_buttons := make(chan variabletypes.ButtonEvent)
    drv_floors  := make(chan int)
    drv_obstr   := make(chan bool)
    drv_stop    := make(chan bool)    
    
    //go elevio.PollButtons(drv_buttons)
    go elevio.PollFloorSensor(drv_floors)
    go elevio.PollObstructionSwitch(drv_obstr)
    go elevio.PollStopButton(drv_stop)
    
    
    for {
        select {
        /*case a := <- drv_buttons:
            fmt.Printf("%+v\n", a)
            elevio.SetButtonLamp(a.Button, a.Floor, true)
          */  
        case a := <- drv_floors:
            //fmt.Printf("Current floor: %+v\n", a)
            if a == config.M_Floors-1 {
                d = variabletypes.MD_Down
            } else if a == 0 {
                d = variabletypes.MD_Up
            }
            elevio.SetMotorDirection(d)
            
            
        case a := <- drv_obstr:
            fmt.Printf("%+v\n", a)
            if a {
                elevio.SetMotorDirection(variabletypes.MD_Stop)
            } else {
                elevio.SetMotorDirection(d)
            }
            
        case a := <- drv_stop:
            fmt.Printf("%+v\n", a)
            for f := 0; f < config.M_Floors; f++ {
                for b := variabletypes.ButtonType(0); b < 3; b++ {
                    elevio.SetButtonLamp(b, f, false)
                }
            }
        }
    }    
}
