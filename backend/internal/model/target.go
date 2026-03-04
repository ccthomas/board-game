package model

// TargetingPattern defines how an ability selects its targets.
type TargetingPattern string

// TargetAlignment defines whether an ability targets enemies, allies, or self.
type TargetAlignment string

const (
	TargetSingle   TargetingPattern = "single"
	TargetAdjacent TargetingPattern = "adjacent" // all tiles directly surrounding the creature
	TargetLine     TargetingPattern = "line"
	TargetCone     TargetingPattern = "cone"
	TargetSelf     TargetingPattern = "self"
	TargetAllyAll  TargetingPattern = "ally_all"
	TargetEnemyAll TargetingPattern = "enemy_all"
)

const (
	AlignEnemy TargetAlignment = "enemy"
	AlignAlly  TargetAlignment = "ally"
	AlignSelf  TargetAlignment = "self"
	AlignAny   TargetAlignment = "any"
)
