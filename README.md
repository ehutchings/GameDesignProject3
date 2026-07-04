Edan Hutchings - Game Design Final Project

How to Play:

Before the game starts, enter the base cost of your towers and your name.  
Any non-number inputs default to a cost of 0, and an empty name will not display above the player's base.  
Use the slider to adjust the difficulty. Towers do normal damage on hard and twice the damage on easy.  
Click the start button when you are ready to play.

The player's base is the blue crystal structure. Typically on the right side of the map.  
The enemy's spawn is the black and red portal structure. Typically on the left side of the map.

Use keys to toggle between each type of tower.  
X - Crossbow: Single target damage, medium fire rate, no special effects. Base cost  
C - Void Launcher: AoE damage by shooting projectile, slow fire rate, stuns enemies on hit. Base cost * 6  
V - Infernal Eye: Single target damage, fast fire rate, long range, burns enemies on hit. Base cost * 2  
B - Snowflake: AoE damage by proximity to tower, medium fire rate, slows enemies on hit. Base cost * 4

Click on an empty tile to place a tower. Towers will not place if it will stop enemies from being able  
to reach the base. If the base price is greater than 0, towers will also cost gold to place.

Water and wall tiles cannot be built on or traversed by enemies.

Enemies drop gold on kill, some enemies drop more gold than others.  
There are regular enemies and fast enemies. Fast enemies have less health and move faster,  
regular enemies move at an average speed and have more health.  
Re-colored versions of enemies are stronger variants with more health.

When an enemy reaches the base, its remaining health will be subtracted from the base's health.  
If the base's health reaches 0 before the player's towers defeat every enemy, the player loses.  
Upon clearing every wave of a given stage, you will enter the next stage.  
Going to a new stage resets your towers and gold. There are 3 reachable stages.  
You are given a grace period between each stage to observe the layout of the new map and place new towers.  
Strategically use towers to extend the enemy's path to win each stage.  
The game will end after clearing all waves in the last stage. You can then close the game.

Things I wanted to finish but did not:  

Layered tile maps - I was not able to make layered tile maps work, but I was still able to add the   
functionality of certain tiles blocking enemies and tower placements.
  
Realistic projectiles -  I am not familiar  with the math to do these things. I used a few shortcuts 
to make the projectiles aim somewhat towards the enemy's  
position rather than using realistic shooting and collisions.

An interface showing the list of possible towers, which one is selected, and each tower's cost.

A way to highlight the current path that enemies will take by holding a certain key.


