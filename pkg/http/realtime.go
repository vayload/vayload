/*
 * SPDX-License-Identifier: MIT
 *
 * Vayload - Container
 *
 * Copyright (c) 2026 Alex Zweiter
 */

package httpi

import (
	"bufio"
	"fmt"
	"time"

	"github.com/goccy/go-json"
	"github.com/vayload/vayload/internal/vayload"
	"github.com/vayload/vayload/pkg/events"
)

func ServerSentEventHandler(realtime events.RealtimeHub) vayload.HttpHandler {
	return func(req vayload.HttpRequest, res vayload.HttpResponse) error {
		// Headers CORS y SSE
		res.SetHeader("Content-Type", "text/event-stream; charset=utf-8")
		res.SetHeader("Cache-Control", "no-cache, no-transform")
		res.SetHeader("Connection", "keep-alive")
		res.SetHeader("X-Accel-Buffering", "no")
		res.SetHeader("Vary", "Origin")

		userId := req.Auth().UserId.String()
		channel := realtime.Subscribe(userId)
		ctx := req.Context()

		fmt.Printf("Client %s connected to SSE\n", userId)

		return res.Status(200).SetBodyStreamWriter(func(w *bufio.Writer) error {
			defer func() {
				fmt.Printf("Client %s disconnected\n", userId)
				realtime.Unsubscribe(userId)
			}()

			initialMsg := fmt.Sprintf("id: %s\nevent: connection\ndata: Connected to SSE\n\n", userId)
			w.WriteString(initialMsg)

			if err := w.Flush(); err != nil {
				return err
			}

			ticker := time.NewTicker(25 * time.Second)
			defer ticker.Stop()

			for {
				select {
				case <-ctx.Done():
					return nil

				case <-ticker.C:
					pingData := fmt.Sprintf("{\"time\": \"%s\"}", time.Now().Format(time.RFC3339))
					msg := fmt.Sprintf("event: ping\ndata: %s\n\n", pingData)

					if _, err := w.WriteString(msg); err != nil {
						return err
					}
					if err := w.Flush(); err != nil {
						return err
					}

				case data, ok := <-channel:
					if !ok {
						return nil
					}

					fmt.Printf("📤 Sending to %s: %s\n", userId, data.EventName)

					jsonData, err := json.Marshal(data.Data)
					if err != nil {
						fmt.Printf("Error serializing: %v\n", err)
						continue
					}

					msg := fmt.Sprintf("id: %s\nevent: %s\ndata: %s\n\n",
						data.Id, data.EventName, jsonData)

					fmt.Printf("📝 Message: %q\n", msg)
					if _, err := w.WriteString(msg); err != nil {
						return err
					}

					// CRÍTICO: Flush después de cada evento
					if err := w.Flush(); err != nil {
						return err
					}

					fmt.Printf("✅ Flushed to %s\n", userId)
				}
			}
		})
	}
}
