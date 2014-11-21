var uri = 'ws://localhost:8080/ws';

var mod = angular.module('go-web-mpd', ['angular-websocket',])

mod.config(['WebSocketProvider', function(WebSocketProvider) {
	WebSocketProvider
	.prefix('')
	.uri(uri);
}]);

mod.service( 'mpd', [ '$rootScope', 'WebSocket', function( $rootScope, WebSocket ) {
	var service = {
		send: function(type, data) {
			WebSocket.send(JSON.stringify({'type':type, 'data':data}));
		},
		onmessage: function(type, data) {
			$rootScope.$emit(type, data);
		},
	};

	WebSocket.onmessage(function(event) {
		var parsedData = JSON.parse(event.data);
		service.onmessage(parsedData.Type, parsedData.Data);
	});

	WebSocket.onopen(function() {
		service.send('init', 'init');
	});
	return service;
}]);

mod.controller('MainCtrl', [ '$scope', '$rootScope', 'mpd', function($scope, $rootScope, mpd) {
	$scope.messages = [];

	$scope.update = function(packet) {
		WebSocket.send(packet.message);
	};

	$scope.play = function() {
		mpd.send("play", "play");
	};

	$scope.stop = function() {
		mpd.send("stop", "stop");
	};

	$rootScope.$on('play', log);
	$rootScope.$on('stop', log);
	$rootScope.$on('init', function(event, data){
		log(event, JSON.parse(data));
	});

	function log(event, data) {
		$scope.messages.push({time:Date.now(), data:data});
		console.log(Date.now() + "::" + data);
	};
}]);





