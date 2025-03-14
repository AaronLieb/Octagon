package conflicts

type player struct {
	Name string `json:"name"`
	Id   int    `json:"id"`
}

type conflict struct {
	Priority int      `json:"priority"`
	Players  []player `json:"players"`
}
