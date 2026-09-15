package program

func toUpsert(req UpsertRequest) Upsert {
	learn := req.Learn
	if learn == nil {
		learn = []string{}
	}
	return Upsert{
		Title:     req.Title,
		Short:     req.Short,
		Full:      req.Full,
		Learn:     learn,
		Age:       req.Age,
		Fee:       req.Fee,
		Level:     req.Level,
		Bonus:     req.Bonus,
		Accent:    req.Accent,
		ImageURL:  req.ImageURL,
		Published: req.Published,
	}
}

func toResponse(p Program) ProgramResponse {
	learn := p.Learn
	if learn == nil {
		learn = []string{}
	}
	return ProgramResponse{
		ID:        p.ID.String(),
		Slug:      p.Slug,
		Title:     p.Title,
		Short:     p.Short,
		Full:      p.Full,
		Learn:     learn,
		Age:       p.Age,
		Fee:       p.Fee,
		Level:     p.Level,
		Bonus:     p.Bonus,
		Accent:    p.Accent,
		ImageURL:  p.ImageURL,
		OrderNo:   p.OrderNo,
		Published: p.Published,
	}
}

func toResponses(items []Program) []ProgramResponse {
	out := make([]ProgramResponse, 0, len(items))
	for _, p := range items {
		out = append(out, toResponse(p))
	}
	return out
}
