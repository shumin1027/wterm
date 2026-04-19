import { useRef } from "react";
import { Terminal } from "@wterm/react";
import { WebSocketTransport } from "@wterm/dom";
// import "@wterm/dom/src/terminal.css";
import "@wterm/react/css";

export default function App() {
  const wsRef = useRef(null);

  function handleReady(term) {
    const proto = location.protocol === "https:" ? "wss" : "ws";
    const ws = new WebSocketTransport({
      url: `${proto}://${location.host}/ws/terminal`,
      // url: "ws://localhost:3000/ws/terminal",
      reconnect: true,
      maxReconnectDelay: 5000,
      onData: (data) => term.write(data),
      onOpen: () => console.log("connected"),
      onClose: () => console.log("disconnected"),
    });
    wsRef.current = ws;
    ws.connect();
    term.onData = (data) => ws.send(data);
    term.onResize = (cols, rows) =>
      ws.send(JSON.stringify({ type: "resize", cols, rows }));
  }

  return (
    <Terminal
      // wasmUrl="/wterm.wasm"
      theme="solarized-dark"
      autoResize
      cursorBlink
      style={{ width: "100%", height: "100%" }}
      onTitle={(title) => {
        document.title = title;
      }}
      onReady={handleReady}
    />
  );
}
