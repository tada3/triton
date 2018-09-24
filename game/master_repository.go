package game

type GameMasterRepo struct {
	masterMap map[string]*GameMaster
}

func NewGameMasterRepo() *GameMasterRepo {
	repo := new(GameMasterRepo)
	repo.masterMap = make(map[string]*GameMaster, 10)
	return repo
}

func (repo *GameMasterRepo) AddGameMaster(key string, gm *GameMaster) {
	repo.masterMap[key] = gm
}

func (repo *GameMasterRepo) GetGameMaster(key string) *GameMaster {
	gm, _ := repo.masterMap[key]
	return gm
}
