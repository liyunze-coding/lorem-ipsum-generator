const copyButton = document.getElementById("copy-button");
if (copyButton) {
	console.log("copy button exists");
	copyButton.addEventListener("click", function () {
		const text = document.getElementById("result").innerText;
		navigator.clipboard.writeText(text).then(
			function () {
				console.log("Copying to clipboard was successful!");
			},
			function (err) {
				console.error("Could not copy text: ", err);
			}
		);
	});
}
