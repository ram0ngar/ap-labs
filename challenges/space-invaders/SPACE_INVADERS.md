# Space Invaders

## Requirements

Refer to [faiface's pixel repository](https://github.com/faiface/pixel)

You can use this project Makefile to get the repository or run:

```
go get github.com/faiface/pixel
```

Ubuntu/Debial-like distributions will need to install ```libgl1-mesa-dev``` and ```xorg-dev``` packages using ```sudo apt install```.

In the projects folder run

```
go mod init github.com/faiface/pixel-examples/platformer  
go build
go run .
```

After the last line the project should start running.

## How it works

There is two Structs, one for the enemy and other for the player


```
type playerShip struct {
	kind  string
	name  string
	x     float64
	y     float64
	alive bool
	vx    float64
	vy    float64
	in    chan msg
}

type enemyShip struct {
	register     map[string]playerShip
	ripChannel   chan msg
	spawnChannel chan msg
	scoreChannel map[string](*chan msg)
}
```

However, there is also Walls that sustain damage from both the player and the enemy. This three are spawned by sending the values to the spawnChannel. Since there is more than one enemy in the map and the walls consists on several bricks, they are sent through a for loop into the channel.

```
func populate() {
	var spacingPx = 50
	var con = 0
	for r := 0; r < 300/spacingPx; r++ {
		for i := 0; i < 400/spacingPx; i++ {
			if con < numEnemies {
				ms.spawnChannel <- msg{cmd: "Add", val: "Enemy", p: playerShip{x: float64(100 + spacingPx*i), y: float64(200 + spacingPx*r), vx: float64(-100 + rand.Intn(1)), vy: 0}}
			}
			con++
		}
	}
	for f := 0; f < 4; f++ {
		for c := 0; c < 3; c++ {
			for r := 0; r < 5; r++ {
				ms.spawnChannel <- msg{cmd: "Add", val: "Wall", p: playerShip{x: float64(100 + f*100 + c*10), y: float64(60 + r*10), vx: 0, vy: 0}}
			}
		}
	}
	ms.spawnChannel <- msg{cmd: "Add", val: "Ship", p: playerShip{x: 100, y: 10, vx: 0, vy: 0}}
}
```

But each of them need a color to identify which is which. By making use of pixel's imdraw we can set the colors of the objects in the Window. 

```
var imd = imdraw.New(nil)
					imd.Clear()
					for _, s := range m.register {
						if s.alive && s.kind == "Ship" {
							imd.Color = colornames.Green
							imd.Push(pixel.V(float64(s.x), float64(s.y)))
							imd.Circle(4, 2)
						} else if s.alive && s.kind == "Enemy" {
							imd.Color = colornames.Blueviolet
							imd.Push(pixel.V(float64(s.x), float64(s.y)))
							imd.Circle(4, 2)
						} else if s.alive && s.kind == "Wall" {
							imd.Color = colornames.White
							imd.Push(pixel.V(float64(s.x), float64(s.y)))
							imd.Circle(5, 4)
						} else if s.alive && s.kind == "Laser" {
							imd.Color = colornames.Orange
							imd.Push(pixel.V(float64(s.x), float64(s.y)))
							imd.Circle(2, 2)
						} else if s.alive && s.kind == "Bomb" {
							imd.Color = colornames.Red
							imd.Push(pixel.V(float64(s.x), float64(s.y)))
							imd.Circle(3, 2)
						}
					}
```

In this game it is known that the enemy both sends lasers to destroy you and advances forward to invade the city. So what happens when an enemy hits you? You lose a live. And if you hit them? They are destroyed and removed from the Window. 

```
if message.cmd == "CheckCollisions" {
				for _, s1 := range m.register {
					if s1.alive && s1.kind == "Laser" {
						for _, s2 := range m.register {
							if s2.alive && (s2.kind == "Enemy" || s2.kind == "Wall") {
								if (s2.x-s1.x)*(s2.x-s1.x)+(s2.y-s1.y)*(s2.y-s1.y) < 40 {
									*m.scoreChannel[s2.name] <- msg{cmd: "Die"}
									score += 10
									*m.scoreChannel[s1.name] <- msg{cmd: "Die"}
								}
							}
						}
					}
					if s1.alive && s1.kind == "Bomb" {
						for _, s2 := range m.register {
							if s2.alive && s2.kind == "Ship" {
								if (s2.x-s1.x)*(s2.x-s1.x)+(s2.y-s1.y)*(s2.y-s1.y) < 40 {
									*m.scoreChannel[s2.name] <- msg{cmd: "One live less"}
									live--
									if live == 0 {
										*m.scoreChannel[s2.name] <- msg{cmd: "Die"}
									}
									*m.scoreChannel[s1.name] <- msg{cmd: "Die"}
								}
							}
							if s2.alive && s2.kind == "Wall" {
								if (s2.x-s1.x)*(s2.x-s1.x)+(s2.y-s1.y)*(s2.y-s1.y) < 40 {
									*m.scoreChannel[s2.name] <- msg{cmd: "Die"}
									*m.scoreChannel[s1.name] <- msg{cmd: "Die"}
								}
							}
						}
					}
					if s1.alive && s1.kind == "Enemy" {
						for _, s2 := range m.register {
							if s2.alive && s2.kind == "Ship" {
								if (s2.x-s1.x)*(s2.x-s1.x)+(s2.y-s1.y)*(s2.y-s1.y) < 40 {
									*m.scoreChannel[s2.name] <- msg{cmd: "One live less"}
									live--
									if live == 0 {
										*m.scoreChannel[s2.name] <- msg{cmd: "Die"}
									}
									*m.scoreChannel[s1.name] <- msg{cmd: "Die"}
								}
							}
							if s2.alive && s2.kind == "Wall" {
								if (s2.x-s1.x)*(s2.x-s1.x)+(s2.y-s1.y)*(s2.y-s1.y) < 40 {
									*m.scoreChannel[s2.name] <- msg{cmd: "Die"}
									*m.scoreChannel[s1.name] <- msg{cmd: "Die"}
								}
							}
						}
					}
				}
			}
```

As for the player's control pixel's Window allow us to set an action with their key listeners.

```
if message.cmd == "KeyPressed" {
				if window != nil {
					if m.scoreChannel["Ship"] != nil {
						if window.Pressed(pixelgl.KeyLeft) {
							*m.scoreChannel["Ship"] <- msg{cmd: "MoveLeft"}
						} else if window.Pressed(pixelgl.KeyRight) {
							*m.scoreChannel["Ship"] <- msg{cmd: "MoveRight"}
						} else if window.Pressed(pixelgl.KeySpace) {
							*m.scoreChannel["Ship"] <- msg{cmd: "Shoot"}
						} else {
							*m.scoreChannel["Ship"] <- msg{cmd: "Stop"}
						}
					}
				}
			}
```