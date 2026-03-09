/*
 * SPDX-License-Identifier: MIT
 *
 * Vayload - Container
 *
 * Copyright (c) 2026 Alex Zweiter
 */

package utils

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"
)

func getLocalIPv4() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", fmt.Errorf("error al obtener direcciones de interfaz: %w", err)
	}

	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String(), nil
			}
		}
	}
	return "", fmt.Errorf("no se encontró una dirección IPv4 no loopback")
}

func GetClientIP(r *http.Request) string {
	var realIP string

	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		ips := strings.Split(xff, ",")
		if len(ips) > 0 {
			realIP = strings.TrimSpace(ips[0])
		}
	}

	if realIP == "" {
		if host, _, err := net.SplitHostPort(r.RemoteAddr); err == nil {
			realIP = host
		} else {
			realIP = r.RemoteAddr
		}
	}

	if os.Getenv("GO_ENV") == "development" {
		localIPv4, err := getLocalIPv4()
		if err != nil {
			fmt.Printf("Advertencia: No se pudo obtener la IPv4 local en modo desarrollo: %v. Usando IP detectada: %s\n", err, realIP)
			return realIP
		}

		switch realIP {
		case "::1", "127.0.0.1", "::ffff:127.0.0.1":
			return localIPv4
		default:
			return realIP
		}
	}

	return realIP
}
