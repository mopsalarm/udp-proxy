package main

import (
	"flag"
	"github.com/Sirupsen/logrus"
	"net"
)

func main() {
	source := flag.String("source", ":5000", "Source port to listen on")
	target := flag.String("target", "127.0.0.1:5001", "Target address to forward to")
	quiet := flag.Bool("quiet", false, "Do not print info logging.")
	flag.Parse()

	if *quiet {
		logrus.SetLevel(logrus.WarnLevel)
	}

	sourceAddress, err := net.ResolveUDPAddr("udp", *source)
	if err != nil {
		logrus.WithError(err).Fatal("Could not resolve source address:", *source)
		return
	}

	targetAddress, err := net.ResolveUDPAddr("udp", *target)
	if err != nil {
		logrus.WithError(err).Fatal("Could not resolve target address:", *source)
		return
	}

	sourceConnection, err := net.ListenUDP("udp", sourceAddress)
	if err != nil {
		logrus.WithError(err).Fatal("Could not listen on address:", *source)
		return
	}

	defer sourceConnection.Close()

	targetConnection, err := net.DialUDP("udp", nil, targetAddress)
	if err != nil {
		logrus.WithError(err).Fatal("Could not 'connect' to target address:", *target)
		return
	}

	for {
		buffer := make([]byte, 64*1024)
		n, addr, err := sourceConnection.ReadFromUDP(buffer)

		if err != nil {
			logrus.WithError(err).Error("Could not receive a packet")
			continue
		}

		logrus.WithField("addr", addr.String()).WithField("bytes", n).Info("Packet received")
		if _, err := targetConnection.Write(buffer[0:n]); err != nil {
			logrus.WithError(err).Warn("Could not forward packet.")
		}
	}
}
