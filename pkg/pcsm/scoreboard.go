package pcsm

// FillDivisions populates the division field on teams based on
// divisions that the system knows about.
func (s *Scoreboard) FillDivisions() {
	dTable := make(map[string]string, len(s.Divisions))
	for _, d := range s.Divisions {
		dTable[d.ID] = d.Name
	}

	for _, t := range s.Teams {
		t.Division = dTable[t.DivisionID]
	}
}
