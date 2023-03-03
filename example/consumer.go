package playground

import "context"

type MyService struct {
	TheInterface SomeInterface
}

func (receiver MyService) UseInterface(ctx context.Context) error {
	_, err := receiver.TheInterface.Do(ctx, "hello")

	return err
}
