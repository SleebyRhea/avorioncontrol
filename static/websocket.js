function toggleHidden(btn, cls) {
	elements = document.getElementsByClassName(cls);
	console.log(elements);
	for (elm of elements){
		if (elm.style.visibility == "hidden") {
			console.log("Object was hidden")
			elm.style.visibility = "visible"
		} else {
			console.log("Object was not hidden")
			elm.style.visibility = "hidden"
		}
	}
}

function prepareChatItem(elm, text) {
	var reType = new RegExp("^\\[[^\\[\\]]+\\] ")
	var reChat = new RegExp("^<([^\\[\\]]+)> ")
	var prefix = text.match(reType, "")
	var msg = text.replace(reType, "")
	var btn = document.createElement("button")
	
	prefix = prefix[0]

	// console.log(prefix)
	switch (true) {
		case (prefix == "[CHAT] "):
			var chatter = msg.match(reChat)[1]
			chatter = chatter.trim()
		
			if (chatter == "Server") {
				btn.classList.add("c-badge")
				btn.classList.add("c-badge--success")
			} else {
				btn.classList.add("c-badge")
				btn.classList.add("c-badge--info")
			}

			msg = msg.replace(reChat, "")
			btn.innerText = chatter
			break;

		case (prefix == "[WARN] "):
			elm.classList.add("serverlog-warn")
			btn.classList.add("c-badge")
			btn.classList.add("c-badge--warning")
			btn.innerText = "Warning"
			break;

		case (prefix == "[ERROR] "):
			elm.classList.add("serverlog-error")
			btn.classList.add("c-badge")
			btn.classList.add("c-badge--error")
			btn.innerText = "Error"
			break;

		default:
			elm.classList.add("serverlog-info")
			btn.classList.add("c-badge")
			btn.classList.add("c-badge--brand")
			btn.innerText = "Info"
			break;
	}

	// elm.innerText += msg
	var t = document.createElement("span")
	t.innerText = " " + msg
	elm.appendChild(t)
	elm.prepend(btn)
	return elm
}

document.addEventListener('DOMContentLoaded', () => {
	var conn;
	var serverlog = document.getElementById("serverlog-window");

	function appendLog(item) {
		var doScroll = serverlog.scrollTop > serverlog.scrollHeight - serverlog.clientHeight;
		serverlog.appendChild(item);
		if (doScroll) {
			serverlog.scrollTop = serverlog.scrollHeight - serverlog.clientHeight;
		}
	}

	if (window["WebSocket"]) {
		var proto = "ws://";

		if (location.protocol == "https:") {
			proto = "wss://";
		}

		conn = new WebSocket(proto + document.location.host + "/ws");

		conn.onclose = function(e) {
			var item = document.createElement("div");
			item.classList.add("c-card__item");
			item.innerHTML = "<b>Connection closed.</>";
			appendLog(item);
		}

		conn.onmessage = function (e) {
			var messages = e.data.split('\n');

			for (var i = 0; i < messages.length; i++) {
				var item = document.createElement("div");
				item.classList.add("c-card__item");
				item.classList.add("c-card__item");
				appendLog(prepareChatItem(item, messages[i]));
			}
		}
	} else {
		var item = document.createElement("div");
		item.innerHTML = "<b>Your browser does not support websockets.</b>"
		appendLog(item);
	}
})