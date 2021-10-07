package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	ctx, _ := signal.NotifyContext(context.Background(), syscall.SIGINT)
	var IP string
	if len(os.Args) == 1 {
		// This option need to docker-compose networking
		// If you use only cli, you need to add real host ip ex.localhost, 127.0.0.1
		IP = "server"
	} else {
		IP = os.Args[1]
	}
	fmt.Println(IP)

	d := net.Dialer{
		Timeout:   time.Second,
		KeepAlive: time.Minute,
	}
	// network_mode: host у клиента в docker-compose.yml
	// "127.0.0.1:9000", ":9000", "localhost:9000"
	// реальный локальный ip "172.20.48.1:9000"
	// Поддержка IPv6 не доступна в Compose 3, только 2.
	// При использовании дефолтной сети DNS имя по имени названия сервиса
	conn, err := d.DialContext(ctx, "tcp", IP+":9000")

	if err != nil {
		log.Fatal(err)
	}
	log.Println(io.Copy(os.Stdout, conn))
}
