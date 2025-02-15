function successStatusCode(event) {
	return event.detail.xhr.status === 200;
}

function eventSuccessful(event) {
	return event.detail.successful;
}

function allowStatusCodes(event, ...codes) {
	if (codes.indexOf(event.detail.xhr.status) !== -1) {
		event.detail.shouldSwap = true;
		event.detail.isError = false;
	}
}
