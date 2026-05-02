package common

type Modifier interface {
	Apply(dice []*Die)
}
