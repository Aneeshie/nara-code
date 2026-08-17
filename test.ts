const ws = new WebSocket("ws://localhost:8000");

console.log("created websocket");

ws.onopen = () => {
    console.log("connected");

    ws.send(JSON.stringify({
        role: "user",
        content: "list directories in this repo"
    }));
};

ws.onmessage = (event) => {
    console.log("received:", event.data);
};

ws.onerror = (event) => {
    console.log("error:", event);
};

ws.onclose = (event) => {
    console.log("closed:", event.code, event.reason);
};