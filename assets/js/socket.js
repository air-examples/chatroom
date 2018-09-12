var ws = new WebSocket("ws://localhost:3333/socket");

ws.onopen = function(evt) {
	console.log("Connection open ...");
};

ws.onmessage = function(evt) {
	console.log("Received Message: ");
	console.log(evt);
};

ws.onclose = function(evt) {
	console.log("Connection closed.");
};

function SendMsg() {
	let msg = document.getElementById("msg_input").value;
	ws.send(msg);
}
