package utils

import (
	"slices"
	"testing"
)

func TestGetConnectionFromClient(t *testing.T) {
	// Arrange
	var rPort int = 69
	var address string = "69.69.69.69"
	//var testConnectionGetter = TestConnectionGetter{}
	var testConnection TestConnectionGetter

	// Act
	var socket = testConnection.GetConnectionFromClient(rPort, address)

	// Assert
	testClientSocket, success := socket.(TestSocket) // Have to do type assertion to make sure TestSocket & not Socket is returned:
	if !success {
		t.Errorf("Test Client Socket Type Assertion - Got: %v, Expected: TestSocket\n", socket)
	}

	if testClientSocket.Port != rPort {
		t.Errorf("Test Client Socket Port - Got: %v, Expected: %v\n", testClientSocket.Port, rPort)
	}

	if testClientSocket.Address != address {
		t.Errorf("Test Client Socket Address - Got %v, Expected: %v\n", testClientSocket.Address, address)
	}
}

func TestGetConnectionFromListener(t *testing.T) {
	var rPort int = 69
	var address string = "0.0.0.0"
	var testConnection TestConnectionGetter

	var socket = testConnection.GetConnectionFromListener(rPort, address)
	// Type assertion:
	testListenerSocket, success := socket.(TestSocket)
	if !success {
		t.Errorf("Test Listener Socket Type Assertion - Got: %v, Expected: TestSocket\n", socket)
	}

	if testListenerSocket.Port != rPort {
		t.Errorf("Test Listener Socket Port - Got: %v, Expected: %v\n", testListenerSocket.Port, rPort)
	}

	if testListenerSocket.Address != address {
		t.Errorf("Test Listener Socket Address = Got: %v, Expected: %v\n", testListenerSocket.Address, address)
	}
}

func testSocketRead(t *testing.T) {
	testReadByteArr := []byte("tiddies")
	var testReadErr error

	var fakeSocket TestSocket

	readReturn, readErr := fakeSocket.Read()
	if readErr != testReadErr {
		t.Errorf("Test Error readErr - Got: %v, Expected: error\n", readErr)
	}

	if !slices.Equal(readReturn, testReadByteArr) {
		t.Errorf("Test Read readReturn - Got: %v, Expected: []byte\n", readReturn)
	}

}