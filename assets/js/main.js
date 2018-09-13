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
			alert('name repeat, please input anothor name.');
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
	console.log(msg);
	let text = '';
	if (msg.type === 'welcome') {
		text = msg.time + '<br/>' + msg.from + ' join chatroom.';
	} else if (msg.type === 'text') {
		text = msg.from + ' (' + msg.time + '): ' + msg.content;
	}
	document.getElementById('msg-list').innerHTML += '<p>'+text+'</p>';
}

function ConnectWebSocket(msg) {
	ws = new WebSocket("ws://localhost:3333/socket");

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
		showmsg('Connection closed.');
	};

	ws.onerror = function(evt) {
		console.log(evt);
	};
	return ws;
}

function SendMsg() {
	let msg = document.getElementById('msg-input').value;
	if (ws === null) {
		ConnectWebSocket(msg);
		document.getElementById('msg-input').value = '';
		return false;
	}
	ws.send(plain({text:msg}));
	document.getElementById('msg-input').value = '';
}
