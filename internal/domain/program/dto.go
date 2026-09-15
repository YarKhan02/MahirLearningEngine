package program

type UpsertRequest struct {
	Title     string   `json:"title"`
	Short     string   `json:"short"`
	Full      string   `json:"full"`
	Learn     []string `json:"learn"`
	Age       string   `json:"age"`
	Fee       string   `json:"fee"`
	Level     string   `json:"level"`
	Bonus     string   `json:"bonus"`
	Accent    string   `json:"accent"`
	ImageURL  string   `json:"imageUrl"`
	Published bool     `json:"published"`
}

type ReorderRequest struct {
	IDs []string `json:"ids"`
}

// ProgramResponse is the full admin-facing shape.
type ProgramResponse struct {
	ID        string   `json:"id"`
	Slug      string   `json:"slug"`
	Title     string   `json:"title"`
	Short     string   `json:"short"`
	Full      string   `json:"full"`
	Learn     []string `json:"learn"`
	Age       string   `json:"age"`
	Fee       string   `json:"fee"`
	Level     string   `json:"level"`
	Bonus     string   `json:"bonus"`
	Accent    string   `json:"accent"`
	ImageURL  string   `json:"imageUrl"`
	OrderNo   int      `json:"orderNo"`
	Published bool     `json:"published"`
}
