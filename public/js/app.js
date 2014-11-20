var uri = 'ws://localhost:8080/ws';

angular.module('go-web-mpd', ['angular-websocket', 'controllers',]).config(function(WebSocketProvider) {
	WebSocketProvider
	.prefix('')
	.uri(uri);
});

angular.module('controllers', []).controller('MainCtrl', function($scope, WebSocket) {
	$scope.messages = [];
	WebSocket.onopen(function() {
		console.log('connection');
		WebSocket.send('message')
	});

	WebSocket.onmessage(function(event) {
		console.log('message: ', event.data);
		$scope.messages.push(event.data);
	});

	$scope.update = function(packet) {
		WebSocket.send(packet.message);
	};
});





