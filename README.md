# GoTraining

Run App: go run .cmd\CliApp
Exit TaskApp: CTRL + C

Run client: go run .cmd\client
Exit TaskApp: CTRL + C

Example message request:
curl -X POST http://localhost:8080/messages   -H "Content-Type: application/json"   -d '{"message":"(message string)","userID":(userID int)}'

Example websocket request without client using browser javascript(dev tools console):
const ws = new WebSocket("ws://localhost:8080/ws-last10");

ws.onmessage = (e) => console.log("msg:", e.data);
ws.onclose = () => console.log("closed");
ws.onerror = (e) => console.error("error", e);

Endpoints:
http://localhost:8080/messages
http://localhost:8080/about/
http://localhost:8080/list/
ws://localhost:8080/ws-last10