package main

import (
	"fmt"
	"io"
	"net"
	"os"
)
import "github.com/spf13/cobra"

func main() {
	// add flags to rootCmd
	rootCmd.Flags().StringP("listenAddress", "l", "0.0.0.0:8080", "address to listen on")
	rootCmd.Flags().StringP("remoteAddress", "r", "localhost:8081", "address to forward to")

	// execute rootCmd with handled error
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// create rootCmd
var rootCmd = &cobra.Command{
	Use:   "port-forward",
	Short: "forwards all connections from listenAddress to remoteAddress",
	Run: func(cmd *cobra.Command, args []string) {
		listenAddress, _ := cmd.Flags().GetString("listenAddress")
		remoteAddress, _ := cmd.Flags().GetString("remoteAddress")
		fmt.Printf("Forwarding connections from %s to %s\n", listenAddress, remoteAddress)
		startForwarder(listenAddress, remoteAddress)
	},
}

// startForwarder starts a forwarder that forwards all connections from listenAddress to remoteAddress
func startForwarder(listenAddress, remoteAddress string) {
	// create listener
	listener, err := net.Listen("tcp", listenAddress)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer listener.Close()

	// accept connections
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}

		// forward connection
		go forwardConnection(conn, remoteAddress)
	}
}

func forwardConnection(conn net.Conn, address string) {
	// connect to remote address
	remoteConn, err := net.Dial("tcp", address)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer remoteConn.Close()

	// create channels to signal when copying is done
	done := make(chan struct{})

	// copy data from local connection to remote connection
	go func() {
		_, err := io.Copy(remoteConn, conn)
		if err != nil {
			fmt.Println(err)
		}
		done <- struct{}{}
	}()

	// copy data from remote connection to local connection
	go func() {
		_, err := io.Copy(conn, remoteConn)
		if err != nil {
			fmt.Println(err)
		}
		done <- struct{}{}
	}()

	// wait for both copying goroutines to finish
	<-done
	<-done
}
