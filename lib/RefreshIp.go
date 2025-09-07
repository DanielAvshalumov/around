package lib

import (
	"fmt"
	"net"
	"time"
)

func RenewIp() error {
	conn, err := net.Dial("tcp", "127.0.0.1:9051")
	if err != nil {
		return fmt.Errorf("failed to connect to tor control: %v", err)
	}
	defer conn.Close()

	pass := "yourmomdoesntbelievemewhenisay"
	auto := fmt.Sprintf("Authenticate  \"%s\"\r\n", pass)
	_, err = conn.Write([]byte(auto))
	if err != nil {
		return fmt.Errorf("Failed to auth %v", err)
	}

	time.Sleep(time.Millisecond * 100)

	_, err = conn.Write([]byte("SIGNAL NEWNYM\r\n"))
	if err != nil {
		return err
	}

	fmt.Println("New Ip Assigned")
	return nil
}
