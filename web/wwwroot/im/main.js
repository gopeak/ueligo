 

var debug = false;
 
var authWindow;

var app;
 
var initApp = function() {
	if (app!=null) { return; }
	app = new App( document.getElementById('canvas'));
  
}

var forceInit = function() {
	initApp()
	document.getElementById('unsupported-browser').style.display = "none";
	return false;
}


// 开始启动
if(Modernizr.canvas && Modernizr.websockets) {
	initApp();
} else {
	document.getElementById('unsupported-browser').style.display = "block";	
	document.getElementById('force-init-button').addEventListener('click', forceInit, false);
}
 
$(function() {
	$('a[rel=external]').click(function(e) {
		e.preventDefault();
		window.open($(this).attr('href'));
	});
});

document.body.onselectstart = function() { return false; }
