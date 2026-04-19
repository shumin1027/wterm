package main

import (
	"embed"
	"encoding/json"
	"log"
	"os"
	"os/exec"

	"github.com/creack/pty"
	"github.com/gofiber/contrib/v3/websocket"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/static"
)

type resizeMsg struct {
	Type string `json:"type"`
	Cols uint16 `json:"cols"`
	Rows uint16 `json:"rows"`
}

func terminalHandler(c *websocket.Conn) {
	shell := os.Getenv("SHELL")
	if shell == "" {
		shell = "/bin/bash"
	}

	cmd := exec.Command(shell, "-l")
	cmd.Env = append(os.Environ(), "TERM=xterm-256color")

	ptmx, err := pty.Start(cmd)
	if err != nil {
		log.Printf("pty start error: %v", err)
		return
	}
	defer func() {
		ptmx.Close()
		cmd.Wait()
	}()

	// PTY output -> WebSocket (binary)
	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := ptmx.Read(buf)
			if err != nil {
				return
			}
			if err := c.WriteMessage(websocket.BinaryMessage, buf[:n]); err != nil {
				return
			}
		}
	}()

	// WebSocket -> PTY
	for {
		msgType, data, err := c.ReadMessage()
		if err != nil {
			return
		}

		switch msgType {
		case websocket.TextMessage:
			// resize: {"type":"resize","cols":N,"rows":N}
			var msg resizeMsg
			if err := json.Unmarshal(data, &msg); err == nil && msg.Type == "resize" {
				pty.Setsize(ptmx, &pty.Winsize{
					Cols: msg.Cols,
					Rows: msg.Rows,
				})
			}
		case websocket.BinaryMessage:
			var msg resizeMsg
			if err := json.Unmarshal(data, &msg); err == nil && msg.Type == "resize" {
				pty.Setsize(ptmx, &pty.Winsize{
					Cols: msg.Cols,
					Rows: msg.Rows,
				})
			} else {
				ptmx.Write(data)
			}
		}
	}
}

//go:embed www/*
var webRoot embed.FS

func main() {
	app := fiber.New(fiber.Config{
		AppName:   "Web Terminal",
		Immutable: true,
		GETOnly:   true,
	})
	app.Use(logger.New(logger.Config{
		ForceColors: true,
	}))
	app.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"*"},
	}))

	// WebSocket upgrade middleware
	app.Use("/ws", func(c fiber.Ctx) error {
		// IsWebSocketUpgrade returns true if the client
		// requested upgrade to the WebSocket protocol.
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	// Serve www
	//app.Get("*", static.New("./www"))
	app.Get("*", static.New("www", static.Config{
		FS:     webRoot,
		Browse: true,
	}))

	app.Get("/ws/terminal", websocket.New(terminalHandler))

	log.Fatal(app.Listen(":3000"))
}
