# Pokemon Type Coverage

Given a pokedex text file, and a type chart, also in a text file, it produces a list of all possible type combinations that together offer total coverage over every type.

It finds combinations of types that, when used in a single party, allow the user to score "super effective" hits against every single pokemon type.

Run it with `go run main.go`.

The first 10 combinations are:

 - Fighting Flying Poison Ground Grass Ice Dark
 - Fighting Poison Ground Rock Ghost Grass Fairy
 - Fighting Flying Ground Steel Electric Ice Dark
 - Fighting Flying Poison Ground Electric Ice Dark
 - Fighting Flying Ground Ghost Steel Electric Ice
 - Fighting Flying Poison Ground Ghost Grass Ice
 - Fighting Flying Ground Ghost Steel Grass Ice
 - Fighting Poison Ground Rock Grass Dark Fairy
 - Fighting Flying Poison Ground Ghost Electric Ice
 - Fighting Flying Ground Steel Grass Ice Dark
