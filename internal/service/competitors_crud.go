package service

import (
	"sort"

	"github.com/GitProger/go-telecom-2025/internal/config"
	"github.com/GitProger/go-telecom-2025/internal/model"
)

type CompetitorService struct { // all competitors are independent on each other
	competitors map[int]*model.Competitor
}

func NewCompetitorService() *CompetitorService {
	return &CompetitorService{
		competitors: make(map[int]*model.Competitor),
	}
}

func (cs *CompetitorService) Register(id int, conf *config.Config) {
	// sync.Pool
	cs.competitors[id] = model.NewCompetitor(id, conf)
}

func (cs *CompetitorService) Delete(id int) {
	delete(cs.competitors, id)
}

func (cs *CompetitorService) Get(id int) *model.Competitor {
	if c, ok := cs.competitors[id]; ok {
		return c
	}
	return nil
}

func (cs *CompetitorService) GetAllMap() map[int]*model.Competitor {
	return cs.competitors
}

func (cs *CompetitorService) GetAll() []*model.Competitor {
	competitors := make([]*model.Competitor, 0, len(cs.competitors))
	for _, c := range cs.competitors {
		competitors = append(competitors, c)
	}
	sort.Slice(competitors, func(i, j int) bool {
		t1 := competitors[i].TimeFromPlannedStart()
		t2 := competitors[j].TimeFromPlannedStart()
		return t1 < t2
	})
	return competitors
}
