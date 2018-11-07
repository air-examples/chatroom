function ajax(){ 
	var ajaxData = { 
		type:arguments[0].type || "POST", 
		url:arguments[0].url || "", 
		async:arguments[0].async || "true", 
		data:arguments[0].data || null, 
		dataType:arguments[0].dataType || "json", 
		contentType:arguments[0].contentType || "application/x-www-form-urlencoded", 
		beforeSend:arguments[0].beforeSend || function(){}, 
		success:arguments[0].success || function(){}, 
		error:arguments[0].error || function(){} 
	} 
	ajaxData.beforeSend() 
	var xhr = createxmlHttpRequest();  
	xhr.responseType=ajaxData.dataType; 
	xhr.open(ajaxData.type,ajaxData.url,ajaxData.async);  
	xhr.setRequestHeader("Content-Type",ajaxData.contentType);  
	xhr.send(convertData(ajaxData.data));  
	xhr.onreadystatechange = function() {  
		if (xhr.readyState == 4) {  
			if(xhr.status == 200){ 
				ajaxData.success(xhr.response) 
			}else{ 
				ajaxData.error(xhr) 
			}  
		} 
	}  
} 

function createxmlHttpRequest() {  
	if (window.ActiveXObject) {  
		return new ActiveXObject("Microsoft.XMLHTTP");  
	} else if (window.XMLHttpRequest) {  
		return new XMLHttpRequest();  
	}  
} 

function convertData(data){ 
	if( typeof data === 'object' ){ 
		var convertResult = "" ;  
		for(var c in data){  
			convertResult+= c + "=" + data[c] + "&";  
		}  
		convertResult=convertResult.substring(0,convertResult.length-1) 
		return convertResult; 
	}else{ 
		return data; 
	} 
}

var tips = {};

getConst();

function getConst() {
	ajax({
		type:'GET',
		url:'/api/const',
		success:function(res) {
			console.log(res);
			tips = res.data;
		},
		error:function(res) {
			tips = {
				'join chatroom':'',
				'message can not be empty':'',
				'name repeat, please input anothor name':'',
				'connection closed':''
			};
			Object.keys(tips).forEach(function(k) {
				tips[k] = k;
			});
		}
	});
}

function InputName() {
	let name = document.getElementById("name").value;
	ajax({
		url:'/api/name',
		data:{'name':name},
		success:function(res) {
			console.log(res);
			window.location.href = '/chat';
		},
		error:function(res) {
			alert(tips['name repeat, please input anothor name']);
		}
	});
}

var ws = null;

function plain() {
	let msg =  {
		'type':arguments[0].type || 'text',
		'content':arguments[0].text || ''
	};
	return JSON.stringify(msg);
}

function showmsg(msg) {
	msg = JSON.parse(msg);
	let frame = document.getElementById('msg-list-frame');
	let list = document.getElementById('msg-list');
	let text = '';
	if (msg.type === 'welcome') {
		text = msg.time + '<br/>' + msg.from + ' ' + tips['join chatroom'];
		text = '<p style="text-align:center;">'+text+'</p>';
	} else if (msg.type === 'text') {
		text = msg.from + ' (' + msg.time + '): ' + msg.content;
		text = '<p>'+text+'</p>';
	}
	list.innerHTML += text;
	if (list.clientHeight > frame.clientHeight) {
		frame.scrollTop = list.clientHeight - frame.clientHeight;
	}
}

function ConnectWebSocket(msg) {
	ws = new WebSocket("ws://localhost:2333/socket");

	ws.onopen = function(evt) {
		console.log("Connection open ...");
		ws.send(plain({type:'welcome'}));
		ws.send(plain({text:msg}));
	};

	ws.onmessage = function(evt) {
		console.log("Received Message: ");
		console.log(evt);
		showmsg(evt.data);
	};

	ws.onclose = function(evt) {
		console.log("Connection closed.");
		showmsg(tips['connection closed']);
	};

	ws.onerror = function(evt) {
		console.log(evt);
	};
	return ws;
}

function SendMsg() {
	setFocus();
	let msg = document.getElementById('msg-input').value;
	if (msg === undefined || msg === null || msg === '') {
		alert(tips['message can not be empty']);
		return false;
	}
	if (ws === null) {
		ConnectWebSocket(msg);
		document.getElementById('msg-input').value = '';
		return false;
	}
	ws.send(plain({text:msg}));
	document.getElementById('msg-input').value = '';
}

function setFocus() {
	document.getElementById('msg-input').focus();
}

setFocus();
