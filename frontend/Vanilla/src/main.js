import "@wterm/dom/css";
import { WTerm, WebSocketTransport } from "@wterm/dom";

const terminal = document.getElementById("terminal");

// 添加主题
terminal.classList.add("theme-monokai");

const term = new WTerm(terminal, {
  cursorBlink: true,
  autoResize: true,
  onTitle(title) {
    document.title = title;
  },
});

function handleReady(term) {
  const proto = location.protocol === "https:" ? "wss" : "ws";
  const ws = new WebSocketTransport({
    url: `${proto}://${location.host}/ws/terminal`,
    // url: "ws://localhost:3000/ws/terminal",
    reconnect: true,
    maxReconnectDelay: 5000,
    // 收到数据
    onData: (data) => term.write(data),
    // 打开连接
    onOpen: () => console.log("connected"),
    // 关闭连接
    onClose: () => console.log("disconnected"),
  });
  ws.connect();
  // 发送数据
  term.onData = (data) => ws.send(data);
  // 终端尺寸自适应
  term.onResize = (cols, rows) =>
    ws.send(JSON.stringify({ type: "resize", cols, rows }));
}

// 初始化控制台
term.init().then((res) => {
  handleReady(term);
});
