package main

//CardData is a simple struct to store card params in DB
type CardData struct {
	power  int
	damage int
}

//CardDB is card database
var CardDB = [4][7]CardData{
	{{3, 2}, {3, 1}, {3, 3}, {2, 2}, {4, 2}, {2, 1}, {4, 3}},
	{{4, 3}, {4, 2}, {4, 4}, {3, 3}, {5, 3}, {3, 2}, {5, 4}},
	{{6, 5}, {6, 4}, {6, 6}, {5, 5}, {7, 5}, {5, 4}, {7, 6}},
	{{8, 7}, {8, 6}, {8, 8}, {7, 7}, {9, 7}, {7, 6}, {9, 8}},
}

/* First var of CardDB
var CardDB = [4][5]CardData{
	{{3, 3}, {4, 3}, {3, 4}, {2, 3}, {3, 2}},
	{{5, 5}, {6, 5}, {5, 6}, {4, 5}, {5, 4}},
	{{7, 7}, {7, 6}, {6, 7}, {8, 7}, {7, 8}},
	{{9, 9}, {9, 8}, {8, 9}, {9, 10}, {10, 9}},
}
*/
