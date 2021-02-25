var events = new EventSource("/events");
events.addEventListener("otp", function(e) {
	var otp = JSON.parse(e.data);
	var code = document.getElementById(otp.id).getElementsByClassName('code')[0];
	var time = document.getElementById(otp.id).getElementsByClassName('time')[0];
	code.innerHTML = otp.code;
	time.value = otp.time;
});
