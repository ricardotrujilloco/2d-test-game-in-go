package main

type vector struct {
	x, y float64
}

type element struct {
	position        vector
	width           float64
	rotation        float64
	active          bool
	boundingCircle  boundingCircle
	logicComponents []logicComponent
	uiComponents    []uiComponent
	attributes      []attribute
}

type updateParameters struct {
	position vector
	rotation float64
	width    float64
	elapsed  float64
}

type gameObject interface {
	isActive() *bool
	getPosition() *vector
	getRotation() *float64
	getWidth() *float64
	update(updateParameters updateParameters) error
	onCollision(otherElement gameObject) error
	draw(parameters drawParameters) error
	getBoundingCircle() boundingCircle
}

type uiComponent interface {
	onDraw(parameters drawParameters) error
}

type logicComponent interface {
	onUpdate(parameters updateParameters) error
}

type attribute interface {
}

var elements []gameObject
