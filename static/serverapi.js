'use strict'

function DOMLoaded() {
	return false
}

function resetElement(e) {
	e.value = null
}

function getElementInsideContainer(pID, chID) {
	var elm = document.getElementById(chID);
	var parent = elm ? elm.parentNode : {};
	return (parent.id && parent.id === pID) ? elm : {};
}

var DEBUG = true
var APIBASE = "/api/"
var APIRE = '\\/api\\/'

var ajaxFullstatus = DOMLoaded
var playerKick     = DOMLoaded
var playerBan      = DOMLoaded
var serverSay      = DOMLoaded
var serverStop     = DOMLoaded
var serverMOTD     = DOMLoaded
var serverTime     = DOMLoaded
var serverStart    = DOMLoaded
var serverStatus   = DOMLoaded
var serverSettle   = DOMLoaded
var serverPassword = DOMLoaded
var serverRestart  = DOMLoaded
var verifyMessage  = DOMLoaded
var getRequester   = DOMLoaded

var scopes = new Map()

class TerraControlAPI {
	constructor(scope, obj) {
		TerraControlAPI.RegisterEndpoint(scope, obj, this)
		this.scope = scope
		this.obj = obj
	}

	static RequestBuilder(s, o, ...args) {
		var r = APIBASE + s + "/" + o + "/";
		for (var v of args) {
			if (typeof v === 'string') {
				r = r + v
			}
		}
		return r;
	}

	static RegisterEndpoint(s, o, n) {
		var scope = scopes.get(s);
		if ( scope == undefined) {
			console.log("TerraControlAPI: RegisterEndpoint: Invalid scope: "+s);
			return false;
		} else {
			scope.set(o, n)
		}
	}

	static Requester(r) {
		var re = new RegExp(APIRE+"[^\\/]+\\/[^\\/]+\\/")
		var s = r.match(re, "g")[0].split("/")
		return scopes.get(s[2]).get(s[3])
	}

	getdata() {
		return false
	}

	onprecall() {
		// Block all calls until the DOM is loaded or this function is overridden
		return DOMLoaded()
	}

	oncomplete(xhttp) {
		if (DEBUG) {
			console.log("DEBUG: TerraControl API: Unimplemented: oncomplete: "+this.request)
		}
	}

	onsuccess(xhttp) {
		if (DEBUG) {
			console.log("DEBUG: TerraControl API: Unimplemented: onsuccess: "+this.request)
		}
	}

	onredirect(xhttp) {
		if (DEBUG) {
			console.log("DEBUG: TerraControl API: Unimplemented: onredirect: "+this.request)
		}
	}

	onfailure(xhttp) {
		if (DEBUG) {
			console.log("DEBUG: TerraControl API: Unimplemented: onfail: "+this.request)
		}
	}

	onbadrequest(xhttp) {
		if (DEBUG) {
			console.log("DEBUG: TerraControl API: Unimplemented: onservererror: "+this.request)
		}
	}

	call(data) {
		if (this.request) {
			this.lastrequest = this.request
		}

		this.request = TerraControlAPI.RequestBuilder(this.scope, this.obj,
			this.getdata(), data);

		if (DEBUG) {
			console.log("Making request: "+this.request)
		}

		var xhttp = new XMLHttpRequest();

		// Confirm that the request is even valid
		if (this.onprecall() === true) {
			xhttp.onreadystatechange = function() {
				if (xhttp.readyState == 4) {
					var r = TerraControlAPI.Requester(this.responseURL)
					r.oncomplete(this)
					switch (true) {
						case (this.status <= 299 && this.status >= 200):
							r.onsuccess(this);
							break;
						case (this.status <= 399 && this.status >= 300):
							r.onredirect(this);
							break;
						case (this.status <= 499 && this.status >= 400):
							r.onfailure(this);
							break;
						case (this.status <= 599 && this.status >= 500):
							console.log("Server Error for API call: ", this)
							r.onservererror(this);
							break;
						default:
							console.log("TerraControl API: Invalid Response: "+this.status)
					}
				}
			} 
			
			// Make the request
			xhttp.open("GET", this.request, true);
			xhttp.send();
			return xhttp.response;
		}
	}
}

var scopes = new Map()

// BEGIN
document.addEventListener('DOMContentLoaded', () => {
	// Only permit the creation of endpoints once the DOM is loaded
	scopes.set("server", new Map())
	scopes.set("player", new Map())
	scopes.set("ajax", new Map())

	ajaxFullstatus = new TerraControlAPI("ajax", "fullstatus")
	playerKick     = new TerraControlAPI("player", "kick")
	playerBan      = new TerraControlAPI("player", "ban")
	serverSay      = new TerraControlAPI("server", "say")
	serverStop     = new TerraControlAPI("server", "stop")
	serverMOTD     = new TerraControlAPI("server", "motd")
	serverTime     = new TerraControlAPI("server", "time")
	serverStart    = new TerraControlAPI("server", "start")
	serverStatus   = new TerraControlAPI("server", "status")
	serverSettle   = new TerraControlAPI("server", "settle")
	serverRestart  = new TerraControlAPI("server", "restart")
	serverPassword = new TerraControlAPI("server", "password")

	// serverSay
	serverSay.onprecall = function() {
		var d = getElementInsideContainer("send-server-message",
			"send-server-message-input");
		console.log(d)
		if (d.classList.contains("c-field--success")) {
			return true
		} else {
			return false
		}
	}

	serverSay.getdata = function() {
		return getElementInsideContainer("send-server-message",
			"send-server-message-input").value;
	}

	serverSay.onsuccess = function() {
		var d = getElementInsideContainer("send-server-message",
			"send-server-message-input");
		resetElement(d)
	}


	// serverMOTD
	serverMOTD.getdata = function() {
		return getElementInsideContainer("send-server-motd",
		"send-server-motd-input").value;
	}

	serverMOTD.onsuccess = function() {
		var d = getElementInsideContainer("send-server-motd",
			"send-server-motd-input");
		resetElement(d); 
	}


	// serverPassword
	serverPassword.getdata = function() {
		return getElementInsideContainer("send-server-password",
		"send-server-password-input").value;
	}

	serverPassword.onsuccess = function() {
		var d = getElementInsideContainer("send-server-password",
			"send-server-password-input");
		resetElement(d); 
	}
	
	serverRestart.onprecall = function() {
		var d = document.getElementById("server-restart-button");
		if (d.classList.contains("c-badge--error")) {
			console.log("Server is currently restarting")
			return false;
		} else {
			d.classList.add("c-badge--error");
			return true;
		}
	}

	serverRestart.oncomplete = function() {
		var d = document.getElementById("server-restart-button");
		if (d.classList.contains("c-badge--error")) {
			d.classList.remove("c-badge--error");
			return true;
		}
	}
	
	// ajaxFullstatus
	ajaxFullstatus.onsuccess = function(xhttp) {
		for (const [key, value] of Object.entries(JSON.parse(xhttp.response))) {
			switch (key) {
				case "WorldName":
					break;

				case "Online":
					break;

				case "Seed":
					document.getElementById("world-seed").innerText =
						"World Seed: " + value
					break;

				case "MOTD":
					document.getElementById("game-motd").innerText =
						"Message of the Day: " + value
					break;

				case "Password":
					document.getElementById("game-password").innerText =
						"Password: " + value
					break;
					
				case "Players":
					var plist = document.getElementById("player-list")

					while (plist.lastElementChild.id != "player-count") {
						plist.removeChild(plist.lastElementChild)
					}

					for (const [_, p] of Object.entries(value)) {
						var name = p.Name
						var ip = p.IP

						console.log("Operating on player: "+name + ":"+ip)

						var pdiv = document.createElement("div")
						var pinput = document.createElement("input")
						var span = document.createElement("span")
						var ipb = document.createElement("button")
						var kick = document.createElement("button")
						var ban = document.createElement("button")

						// Primary container class
						pdiv.classList.add("c-card__item")
						pdiv.classList.add("c-input-group")
						pdiv.classList.add("player-container")

						// The input here is the first child class
						pinput.classList.add("c-field")
						pinput.setAttribute("value", name)
						pinput.readOnly = true
						
						// Span is the second, which contains our buttons
						span.classList.add("c-input-group")

						// Our Buttons
						for (var elm of [ipb, kick, ban]) {
							elm.classList.add("c-input-group")
							elm.classList.add("c-button")
							elm.setAttribute("type", "button")
							elm.value = name
						}

						ipb.classList.add("c-button--brand")
						kick.classList.add("c-button--warning")
						ban.classList.add("c-button--error")
						
						ipb.innerText = ip
						kick.innerText = 'Kick'
						ban.innerText = 'Ban'

						kick.addEventListener('click', function() {
							playerKick.call(this.value)
						})

						ban.addEventListener('click', function() {
							playerBan.call(this.value)
						})

						span.append(ipb, kick, ban)
						pdiv.append(pinput, span)
						plist.append(pdiv)
					}
					break;

				case "PlayerCount":
					document.getElementById("player-count").innerText = 
						"Players: " + value
					break;

				case "Loglevel":
					break;
					
				case "Version":
					break;
			}
		}
	}

	verifyMessage = function (elm, min, max) {
		var i = getElementInsideContainer("send-server-div",
			"send-server-message-button")
	
		if (i.classList.contains("c-button--brand")) {
			i.classList.remove("c-button--brand")
		}
	
		if (elm.value.length > max || elm.value.length < min) {
			elm.classList.remove("c-field--success")
			elm.classList.add("c-field--error")
			i.classList.add("c-button--error")
			if (i.classList.contains("c-button--success")) {
				i.classList.remove("c-button--success")
			}
		} else {
			elm.classList.remove("c-field--error")
			elm.classList.add("c-field--success")
			i.classList.add("c-button--success")
			if (i.classList.contains("c-button--error")) {
				i.classList.remove("c-button--error")
			}
		}
	}

	// playerKick
	playerKick.oncomplete = function() {
		setTimeout(function() { ajaxFullstatus.call() }, 3000)
	}

	setInterval(function(){ ajaxFullstatus.call() }, 10 * 1000)

	if (DEBUG) {
		console.log("DOM is ready, and javascript is loaded.")
	}

	DOMLoaded = function() {
		return true
	}
})