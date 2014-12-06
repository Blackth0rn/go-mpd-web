var uri = 'ws://localhost:8080/ws';

var mod = angular.module('go-web-mpd', ['angular-websocket', 'ui.slider',])

mod.config(['WebSocketProvider', function(WebSocketProvider) {
	WebSocketProvider
	.prefix('')
	.uri(uri);
}]);

mod.service( 'mpd', [ '$rootScope', 'WebSocket', function( $rootScope, WebSocket ) {
	var service = {
		send: function(type, data, broadcast) {
			WebSocket.send(JSON.stringify({'cmd':type, 'data':data, 'token':this.token}));
		},
		onmessage: function(type, data) {
			$rootScope.$emit(type, data);
		},
		token: -1,
	};

	WebSocket.onmessage(function(event) {
		var parsedData = JSON.parse(event.data);
		if( parsedData.Cmd == 'register' )
		{
			service.token = parsedData.Token;
			service.send('init', 'init');
		}
		else
			service.onmessage(parsedData.Cmd, parsedData);
	});

	WebSocket.onopen(function() {
	});
	return service;
}]);

mod.controller('MainCtrl', [ '$scope', '$rootScope', 'mpd', function($scope, $rootScope, mpd) {
	$scope.messages = [];
	$scope.mpd_state = {};
	$scope.state = {};

	$scope.update = function(packet) {
		WebSocket.send(packet.message);
	};

	$scope.play = function() {
		mpd.send("play", "play");
	};

	$scope.pause = function() {
		mpd.send("pause", "pause");
	}

	$scope.stop = function() {
		mpd.send("stop", "stop");
	};

	$scope.setVolume = function() {
		mpd.send("setVolume", $scope.state.volume);
	}

	$rootScope.$on('play', playPauseStop);
	$rootScope.$on('pause', playPauseStop);
	$rootScope.$on('stop', playPauseStop);
	$rootScope.$on('init', init);
	$rootScope.$on('setVolume', init);

	function log(event, data) {
		$scope.messages.push({time:Date.now(), data:data});
		console.log(Date.now() + "::" + data);
	};

	function init(event, data) {
		log(event, data);
		$scope.mpd_state = data.Attr;
		$scope.state.showPlay = updatePlayPause();
	}

	function playPauseStop(event, data) {
		$scope.mpd.state = data.Attr;
		$scope.state.showPlay = updatePlayPause();
	}

	function updatePlayPause() {
		if ( $scope.mpd_state.state = 'play' )
		{
			return 'pause';
		}
		else
		{
			return 'play';
		}
	}

}]);





