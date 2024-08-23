var dayOfWeek = new Date().toLocaleString('en-us', {  weekday: 'long' })
document.getElementsByTagName("footer")[0].getElementsByTagName("div")[0].innerHTML = "Have a nice " + dayOfWeek; //change content

var navLinks = document.querySelectorAll("nav a");
for (var i = 0; i < navLinks.length; i++) {
	var link = navLinks[i]
	if (link.getAttribute('href') == window.location.pathname) {
		link.classList.add("live");
		break;
	}
}