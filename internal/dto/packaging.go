package dto

import (
	"workshop-1/internal/domain/strategy"
	"workshop-1/pkg/pvz/v1"
)

func ToPackaging(req pvz.PackagingType) strategy.Packaging {
	switch req {
	case pvz.PackagingType_WRAP:
		return strategy.Wrap{}
	case pvz.PackagingType_BAG:
		return strategy.Bag{}
	case pvz.PackagingType_BOX:
		return strategy.Box{}
	default:
		return nil
	}
}
