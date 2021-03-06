package main

type bulletMover struct {
	position vector
	speed    float64
	state    ElementState
}

func newBulletMover(speed float64) *bulletMover {
	return &bulletMover{
		speed: speed,
	}
}

func (mover *bulletMover) onUpdate(parameters updateParameters) error {

	if mover.position.x == 0 && mover.position.y == 0 {
		mover.position.x = parameters.position.x
		mover.position.y = parameters.position.y
		mover.state = Active
	}

	// mover.position.x += bulletSpeed * parameters.elapsed
	mover.position.y -= bulletSpeed * parameters.elapsed

	if mover.position.x > float64(screenWidth) || mover.position.x < float64(0) ||
		mover.position.y > float64(screenHeight) || mover.position.y < float64(0) {
		mover.position.x = 0
		mover.position.y = 0
		mover.state = Inactive
		return nil
	}

	return nil
}
