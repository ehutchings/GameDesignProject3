package main

const (
	DEFAULT_STAGE_LENGTH = 5
)

type stageWaves struct {
	waves  []*wave
	length int
}

type wave struct {
	enemies       []enemy
	spawnInterval int
	length        int
}

func newWavesForStage(spawnInterval int, stageLength int, waveLength int) *stageWaves {
	newWaves := &stageWaves{
		waves: make([]*wave, stageLength),
	}
	for i := 0; i < stageLength; i++ {
		newWaves.waves[i] = newWave(spawnInterval, waveLength)
	}
	return newWaves
}

func newWave(spawnInterval int, length int) *wave {
	newWave := wave{
		enemies:       make([]enemy, 0),
		spawnInterval: spawnInterval,
	}
	for index := 0; index < length; index++ {
		newEnemy := newEnemy(0, 0, 2)
		newWave.enemies = append(newWave.enemies, *newEnemy)
	}
	return &newWave
}

func (wave *wave) removeEnemyInFront() enemy {
	var currentEnemy *enemy = nil
	if len(wave.enemies) > 1 {
		currentEnemy = &wave.enemies[0]
		wave.enemies = wave.enemies[:1]
	} else {
		currentEnemy = &wave.enemies[0]
		wave.enemies = wave.enemies[:0]
	}
	return *currentEnemy
}

func (stageWaves *stageWaves) removeWaveInFront() wave {
	var currentWave *wave = nil
	if len(stageWaves.waves) > 1 {
		currentWave = stageWaves.waves[0]
		stageWaves.waves = stageWaves.waves[:1]
	} else {
		currentWave = stageWaves.waves[0]
		stageWaves.waves = stageWaves.waves[:0]
	}
	return *currentWave
}

func (stageWaves *stageWaves) getNextWave() *wave {
	if len(stageWaves.waves) != 0 {
		nextWave := stageWaves.removeWaveInFront()
		return &nextWave
	}
	return nil
}
