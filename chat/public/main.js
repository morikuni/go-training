function start() {
	var name = document.getElementById("name")
	name.disabled = true
	document.getElementById("connect").disabled = true
	var url = "ws://" + location.host + "/chat";
	webSocket = new WebSocket(url);
	var messages = document.getElementById("messages")


	webSocket.onopen = function () {
		console.log("connected")
		webSocket.send(JSON.stringify({name: name.value}));
	};

	webSocket.onmessage = function(event) {
		if (event && event.data) {
			var elem = messageElement(JSON.parse(event.data));
			messages.appendChild(elem)
		}
	};

	webSocket.onerror = function (error) {
		console.log('WebSocket Error ' + error);
	};

	webSocket.onclose= function (error) {
		console.log('WebSocket Close ' + error);
	};
}

function messageElement(msg) {
	var messageWrapper = document.createElement("div")
	messageWrapper.className = "messageWrapper"
	var sender = document.createElement("span")
	sender.className = "sender"
	sender.innerHTML = msg.sender
	var message = document.createElement("span")
	message.className = "message"
	message.innerHTML = msg.message
	messageWrapper.appendChild(sender)
	messageWrapper.appendChild(message)
	return messageWrapper
}

function sendMessage() {
	json = {sender: null, message: document.getElementById("msg").value}
	webSocket.send(JSON.stringify(json))
}
