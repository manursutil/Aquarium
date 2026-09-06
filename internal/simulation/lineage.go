package simulation

import "sort"

type LineageRecord struct {
	ID, ParentA, ParentB FishID
	BornIn, EndedIn      int
	EndReason            string
	Elite                bool
	Fitness              float32
	Genome               Genome
}

func recordFromFish(f Fish) LineageRecord {
	return LineageRecord{ID: f.ID, ParentA: f.ParentA, ParentB: f.ParentB, BornIn: f.BornIn, Elite: f.Elite, Fitness: fitness(f), Genome: f.Genome}
}

func (a *Simulation) finalizeLineage(f Fish, reason string) {
	r := recordFromFish(f)
	
	r.EndedIn, r.EndReason = a.generation, reason
	if a.lineage == nil {
		a.lineage = make(map[FishID]LineageRecord)
	}
	
	a.lineage[f.ID] = r
}

func (a *Simulation) Lineage(id FishID) (LineageRecord, bool) {
	for _, f := range a.fish {
		if f.ID == id {
			return recordFromFish(f), true
		}
	}
	
	r, ok := a.lineage[id]
	
	return r, ok
}

func (a *Simulation) LineageRecords() []LineageRecord {
	records := make([]LineageRecord, 0, len(a.lineage))
	
	active := make(map[FishID]LineageRecord, len(a.fish))
	
	for _, f := range a.fish {
		active[f.ID] = recordFromFish(f)
	}
	
	for id, r := range a.lineage {
		if current, ok := active[id]; ok {
			r = current
		}
		records = append(records, r)
	}
	
	sort.Slice(records, func(i, j int) bool {
		if records[i].BornIn != records[j].BornIn {
			return records[i].BornIn < records[j].BornIn
		}
		return records[i].ID < records[j].ID
	})
	
	return records
}

func (a *Simulation) evictLineage() {
	for _, r := range a.LineageRecords() {
		if r.EndedIn > 0 && r.EndedIn < a.generation-a.config.LineageGenerations {
			delete(a.lineage, r.ID)
		}
	}
}
