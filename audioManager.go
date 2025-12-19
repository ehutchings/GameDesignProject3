package main

import "github.com/hajimehoshi/ebiten/v2/audio"

type audioManager struct {
	regularEnemySoundPlayer *audio.Player
	fastEnemySoundPlayer    *audio.Player
	voidLauncherSoundPlayer *audio.Player
	infernalEyeSoundPlayer  *audio.Player
	crossbowSoundPlayer     *audio.Player
	snowflakeSoundPlayer    *audio.Player
	victorySoundPlayer      *audio.Player
	defeatSoundPlayer       *audio.Player
}

func (audManager *audioManager) playVoidLauncherSound() {
	audManager.voidLauncherSoundPlayer.Rewind()
	audManager.voidLauncherSoundPlayer.Play()
}

func (audManager *audioManager) playCrossbowSound() {
	audManager.crossbowSoundPlayer.Rewind()
	audManager.crossbowSoundPlayer.Play()
}

func (audManager *audioManager) playInfernalEyeSound() {
	audManager.infernalEyeSoundPlayer.Rewind()
	audManager.infernalEyeSoundPlayer.Play()
}

func (audManager *audioManager) playSnowflakeSound() {
	audManager.snowflakeSoundPlayer.Rewind()
	audManager.snowflakeSoundPlayer.Play()
}
func (audManager *audioManager) playRegularEnemyDeathSound() {
	audManager.regularEnemySoundPlayer.Rewind()
	audManager.regularEnemySoundPlayer.Play()
}

func (audManager *audioManager) playFastEnemyDeathSound() {
	audManager.fastEnemySoundPlayer.Rewind()
	audManager.fastEnemySoundPlayer.Play()
}

func (audManager *audioManager) playVictorySound() {
	audManager.victorySoundPlayer.Rewind()
	audManager.victorySoundPlayer.Play()
}

func (audManager *audioManager) playDefeatSound() {
	audManager.defeatSoundPlayer.Rewind()
	audManager.defeatSoundPlayer.Play()
}
