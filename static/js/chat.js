function toggleDarkMode() {
	document.body.classList.toggle("dark-mode");
	document.getElementById("chat-container").classList.toggle("dark-mode");
	document.getElementById("topbar").classList.toggle("dark-mode");
	document.getElementById("sun-icon").classList.toggle("dark-mode");
	document.getElementById("moon-icon").classList.toggle("dark-mode");
	document.getElementById("chat-window").classList.toggle("dark-mode");
	document.getElementById("input-area").classList.toggle("dark-mode");
	const messages = document.querySelectorAll(".message");
	messages.forEach(message => {
		message.classList.toggle("dark-mode");
	});
}

document.addEventListener("DOMContentLoaded", function() {
	document.body.classList.toggle("chat-page");
	document.body.classList.toggle("dark-mode");

	let messages = document.querySelectorAll(".message");

	if (messages.length > 0) {
		let last_message = messages[messages.length - 1];
		last_message.scrollIntoView({ behavior: "smooth", block: "end" });
		let chat_window = document.getElementById("chat-window");
		setTimeout(() => {
			chat_window.scrollTop = chat_window.scrollHeight
		}, 500);
	}

	document.getElementById("dark-mode-toggle").addEventListener("click", toggleDarkMode);

	document.getElementById("message-input").addEventListener("keypress", function(event) {
		if (event.key === "Enter") {
			event.preventDefault();
			document.getElementById("send-button").click();
		}
	});

	document.getElementById("send-button").addEventListener("htmx:beforeRequest", function(event) {
		var input_elem = document.getElementById("message-input");
		var msg = input_elem.value.trim();

		if (msg === "") {
			event.preventDefault();
			return;
		}

		input_elem.value = "";

		var new_div = document.createElement("div");
		new_div.classList.add("message", "user-message", "dark-mode");

		var new_span = document.createElement("span");
		new_span.classList.add("nickname");
		new_span.textContent = "You";
		new_div.appendChild(new_span);

		new_div.appendChild(document.createElement("br"));

		var new_pre = document.createElement("pre");
		new_pre.textContent = msg;
		new_div.appendChild(new_pre);

		var chat_window = document.getElementById("chat-window");
		chat_window.appendChild(new_div);

		new_div.scrollIntoView({ behavior: "smooth", block: "end" });
		setTimeout(() => {
			chat_window.scrollTop = chat_window.scrollHeight
		}, 200);
	});

	document.getElementById("chat-window").addEventListener("htmx:afterSwap", function(event) {
		let messages = document.querySelectorAll(".message");
		let last_message = messages[messages.length - 1];
		last_message.scrollIntoView({ behavior: "smooth", block: "end" });
		var chat_window = document.getElementById("chat-window");
		setTimeout(() => {
			chat_window.scrollTop = chat_window.scrollHeight
		}, 300);
	});
});
