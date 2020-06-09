package main

import (
	"flag"
	"fmt"
	"math/rand"
	"sync/atomic"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font/basicfont"
)

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

type msg struct {
	cmd string
	val string
	p   playerShip
}

var gos uint64

var window *pixelgl.Window

var ms = enemyShip{register: make(map[string]playerShip), scoreChannel: make(map[string](*chan msg)), spawnChannel: make(chan msg, 20000), ripChannel: make(chan msg, 20000)}

var channelSpeed int
var shipSpeed int
var numEnemies int
var score int
var live int

func callGo(f func()) {
	atomic.AddUint64(&gos, 1)
	go f()
}

func (m *enemyShip) conductor() {
	populate()
	callGo(ms.actions)
	var a = 0
	for range time.NewTicker(time.Duration(shipSpeed) * time.Millisecond).C {
		a = a + 1
		if a%2 == 0 {
			ms.spawnChannel <- msg{cmd: "Spawn"}
		}
		if a%2 == 0 {
			ms.spawnChannel <- msg{cmd: "CheckCollisions"}
		}
		if a%2 == 0 {
			ms.spawnChannel <- msg{cmd: "KeyPressed"}
		}
	}
}

func (m *enemyShip) actions() {
	for {
		select {
		case message := <-m.ripChannel:
			if message.cmd == "Set" {
				m.register[message.p.name] = playerShip{name: message.p.name, kind: message.p.kind, x: message.p.x, y: message.p.y, alive: message.p.alive}
			} else if message.cmd == "Remove" {
				delete(m.register, message.p.name)
			}
		case message := <-m.spawnChannel:
			if message.cmd == "Add" {
				name := message.val
				if name != "Ship" {
					name = fmt.Sprintf(message.val, time.Now())
				}
				p := playerShip{name: name, kind: message.val, x: message.p.x, y: message.p.y, alive: true, vx: message.p.vx, vy: message.p.vy, in: make(chan msg)}
				m.scoreChannel[p.name] = &p.in
				callGo(p.conductor)
			}
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
			if message.cmd == "Spawn" {
				if window != nil {
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
					window.Clear(colornames.Black)
					imd.Draw(window)
					basicAtlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)
					basicTxt := text.New(pixel.V(100, 500), basicAtlas)
					fmt.Fprintln(basicTxt, "Score: ", score)
					fmt.Fprintln(basicTxt, "Lives: ", live)
					basicTxt.Draw(window, pixel.IM)
					window.Update()
				}
			}
		}
	}
}

func (p *playerShip) conductor() {
	callGo(p.actions)
	for range time.NewTicker(time.Duration(channelSpeed) * time.Millisecond).C {
		p.in <- msg{cmd: "Move"}
	}
}

func (p *playerShip) actions() {
	for {
		m := <-p.in
		if m.cmd == "Die" {
			p.alive = false
			ms.ripChannel <- msg{cmd: "Remove", p: *p}
			if p.kind == "Ship" {
				fmt.Println("Game Over")
				<-make(chan bool)
			}
		} else if m.cmd == "MoveLeft" {
			if p.x > 20 {
				p.vx = -80
			}
		} else if m.cmd == "Stop" {
			p.vx = 0
		} else if m.cmd == "MoveRight" {
			if p.x < 490 {
				p.vx = 80
			}
		} else if m.cmd == "Shoot" {
			ms.spawnChannel <- msg{cmd: "Add", val: "Laser", p: playerShip{x: p.x, y: p.y, vx: 0, vy: 100}}
		}
		if m.cmd == "Move" {
			var xPixPerBeat = p.vx / 1000 * float64(shipSpeed)
			var yPixPerBeat = p.vy / 1000 * float64(shipSpeed)
			p.x = p.x + xPixPerBeat
			if p.alive && p.kind == "Enemy" {
				if rand.Intn(950) < 1 {
					ms.spawnChannel <- msg{cmd: "Add", val: "Bomb", p: playerShip{x: p.x, y: p.y, vx: 0, vy: -100}}
				}
				if p.x > 500 || p.x < 10 {
					p.vx = -p.vx
					p.y = p.y - 10
				}
			} else {
				p.y = p.y + yPixPerBeat
			}
			if p.kind == "Laser" && p.y > 600 {
				p.alive = false
				ms.ripChannel <- msg{cmd: "Remove", p: *p}
			} else if p.kind == "Bomb" && p.y < 0 {
				p.alive = false
				ms.ripChannel <- msg{cmd: "Remove", p: *p}
			}
			if p.alive {
				ms.ripChannel <- msg{cmd: "Set", p: *p}
			}
		}
	}
}

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

func start() {
	window, _ = pixelgl.NewWindow(pixelgl.WindowConfig{Title: "Space Invaders", Bounds: pixel.R(0, 0, float64(512), float64(512)), VSync: true})
	window.SetPos(window.GetPos().Add(pixel.V(0, 1)))
	callGo(ms.conductor)
	<-make(chan bool)
}

func main() {
	channelSpeed = 8
	shipSpeed = 17
	live = 10
	flag.IntVar(&numEnemies, "Enemies", 50, "Number of Enemies")

	flag.Parse()
	pixelgl.Run(start)
}