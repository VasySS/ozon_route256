package strategy

import "workshop-1/internal/domain"

type Packaging interface {
	Apply(order *domain.Order) error
}

func NewPackaging(name string) (Packaging, error) {
	switch name {
	case "wrap":
		return Wrap{}, nil
	case "bag":
		return Bag{}, nil
	case "box":
		return Box{}, nil
	default:
		return nil, PackagingError{"неизвестный тип упаковки: " + name}
	}
}

type Wrap struct{}

func (f Wrap) Apply(order *domain.Order) error {
	order.Price += 1

	return nil
}

type Bag struct{}

func (b Bag) Apply(order *domain.Order) error {
	if order.Weight >= 10 {
		return PackagingError{"заказ слишком тяжелый для пакета"}
	}

	order.Price += 5

	return nil
}

type Box struct{}

func (b Box) Apply(order *domain.Order) error {
	if order.Weight >= 30 {
		return PackagingError{"заказ слишком тяжелый для коробки"}
	}

	order.Price += 20

	return nil
}

type PackagingWithWrap struct {
	MainPackaging Packaging
}

func (p PackagingWithWrap) Apply(order *domain.Order) error {
	if p.MainPackaging == nil {
		return PackagingError{"не задана основная упаковка"}
	}

	if _, isWrap := p.MainPackaging.(Wrap); isWrap {
		return PackagingError{"нельзя использовать плёнку дважды"}
	}

	wrap := Wrap{}
	if err := wrap.Apply(order); err != nil {
		return err
	}

	if err := p.MainPackaging.Apply(order); err != nil {
		return err
	}

	return nil
}
