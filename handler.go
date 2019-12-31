package main

import (
	logger "github.com/sirupsen/logrus"
)

// Handler interface contains the methods that are required
type Handler interface {
	Init() error
	ObjectCreated(obj interface{})
	ObjectDeleted(obj interface{})
	ObjectUpdated(objOld, objNew interface{})
}

// TestHandler is a sample implementation of Handler
type TestHandler struct{}

// Init handles any handler initialization
func (t *TestHandler) Init() error {
	logger.Info("TestHandler.Init")
	return nil
}

// ObjectCreated is called when an object is created
func (t *TestHandler) ObjectCreated(obj interface{}) {
	logger.Info("TestHandler.ObjectCreated")
}

// ObjectDeleted is called when an object is deleted
func (t *TestHandler) ObjectDeleted(obj interface{}) {
	logger.Info("TestHandler.ObjectDeleted")
}

// ObjectUpdated is called when an object is updated
func (t *TestHandler) ObjectUpdated(objOld, objNew interface{}) {
	logger.Info("TestHandler.ObjectUpdated")
}
