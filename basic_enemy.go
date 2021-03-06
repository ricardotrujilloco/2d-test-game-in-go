package main

import (
	"fmt"
	"github.com/veandco/go-sdl2/sdl"
	"reflect"
)

const (
	basicEnemyWidth  = 105
	basicEnemyHeight = 72
	explosionSpeed   = 0.1
)

type enemy struct {
	element
	state ElementState
}

func (elem *enemy) isActive() bool {
	return elem.state == Active ||
		elem.state == Destroying
}

func (elem *enemy) getPosition() vector {
	return elem.position
}

func (elem *enemy) getRotation() float64 {
	return elem.rotation
}

func (elem *enemy) getWidth() float64 {
	return elem.width
}

func (elem *enemy) update(updateParameters updateParameters) error {
	var err error = nil
	updateParameters.state = elem.state
	for _, comp := range elem.logicComponents {
		err = comp.onUpdate(updateParameters)
	}
	elem.onAnimatorUpdated()
	return err
}

func (elem *enemy) onCollision(otherElement gameObject) error {
	switch otherElement.(type) {
	case *bullet:
		elem.onBulletCollision()
	case *enemy:
		elem.onEnemyCollision()
	}
	return nil
}

func (elem *enemy) draw() error {
	parameters := multiSpriteDrawParameters{
		position: elem.getPosition(),
		rotation: elem.getRotation(),
	}
	circleParameters := circleDrawParameters{
		drawParameters: &parameters,
		radius:         int32(elem.getBoundingCircle().radius),
	}
	for _, comp := range elem.uiComponents {
		var err error = nil
		if reflect.TypeOf(comp) == reflect.TypeOf(&circleRenderer{}) {
			err = comp.onDraw(&circleParameters)
		} else {
			err = comp.onDraw(&parameters)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func (elem *enemy) getBoundingCircle() *boundingCircle {
	return elem.boundingCircle
}

func (elem *enemy) onAnimatorUpdated() {
	if component, ok := elem.logicComponents[Animator]; ok {
		animator := component.(*animator)
		if animator.finished {
			elem.state = Inactive
		}
	}
}

func (elem *enemy) onBulletCollision() {
	isVulnerableToBullets := false
	for _, attr := range elem.attributes {
		switch attr.(type) {
		case *vulnerableToBullets:
			isVulnerableToBullets = true
		}
	}
	if isVulnerableToBullets {
		elem.state = Destroying
		elem.setAnimatorState(Destroying)
	}
}

func (elem *enemy) onEnemyCollision() {
	if elem.state == Active {
		elem.state = Destroying
		elem.setAnimatorState(Destroying)
	}
}

func (elem *enemy) setAnimatorState(state ElementState) {
	if component, ok := elem.logicComponents[Animator]; ok {
		animator := component.(*animator)
		animator.setSequence(state)
	}
}

func (elem *enemy) setJumperState(state ElementState) {
	if component, ok := elem.logicComponents[JumpMover]; ok {
		jumpMover := component.(*jumpMover)
		jumpMover.setState(state)
	}
}

func newBasicEnemy(renderer *sdl.Renderer, position vector) enemy {
	destroyingSampleRate := 15.0
	basicEnemyRadiusScaleFactor := 0.25
	basicEnemyInitialRadius := (810 / 4) * basicEnemyRadiusScaleFactor // From sprite dimensions
	basicEnemyFinalRadius := (810 / 2) * basicEnemyRadiusScaleFactor   // From sprite dimensions
	animator := newAnimator(getEnemySequences(destroyingSampleRate), Active)
	circle := &boundingCircle{center: position, radius: basicEnemyInitialRadius}
	boundingCircles := []*boundingCircle{circle}
	boundingCircleScaler := newBoundingCircleScaler(boundingCircles, basicEnemyFinalRadius)
	return enemy{
		state: Active,
		element: element{
			position: position,
			rotation: 180,
			logicComponents: map[LogicComponentType]logicComponent{
				Animator:             animator,
				BoundingCircleScaler: boundingCircleScaler,
			},
			attributes: []attribute{&vulnerableToBullets{}},
			uiComponents: []uiComponent{
				newMultiSpriteRenderer(
					renderer,
					getEnemyUiSequences(renderer),
					animator,
					basicEnemyRadiusScaleFactor,
				),
				newCircleRenderer(
					renderer,
					boundingCircles,
				),
			},
			boundingCircle: circle,
		},
	}
}

func getEnemySequences(
	destroyingSampleRate float64,
) map[ElementState]*sequence {
	idleSequence, err := newSequence("data/sprites/bomb/idle", 10, true, false)
	if err != nil {
		panic(fmt.Errorf("creating idle sequence: %v", err))
	}
	destroySequence, err := newSequence("data/sprites/bomb/destroy", destroyingSampleRate, false, true)
	if err != nil {
		panic(fmt.Errorf("creating onBulletCollision sequence: %v", err))
	}
	sequences := map[ElementState]*sequence{
		Active:     idleSequence,
		Destroying: destroySequence,
	}
	return sequences
}

func getEnemyUiSequences(renderer *sdl.Renderer) map[ElementState]*multiSpriteRendererSequence {
	idleSequenceUi, err := newMultiSpriteRendererSequence("data/sprites/bomb/idle", renderer)
	if err != nil {
		panic(fmt.Errorf("creating idle sequence: %v", err))
	}
	destroySequenceUi, err := newMultiSpriteRendererSequence("data/sprites/bomb/destroy", renderer)
	if err != nil {
		panic(fmt.Errorf("creating onBulletCollision sequence: %v", err))
	}
	uiSequences := map[ElementState]*multiSpriteRendererSequence{
		Active:     idleSequenceUi,
		Destroying: destroySequenceUi,
	}
	return uiSequences
}
