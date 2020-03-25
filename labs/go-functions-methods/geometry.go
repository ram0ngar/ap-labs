// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 156.

// Package geometry defines simple types for plane geometry.
//!+point
package main

import (
	"fmt"
	"math"
	"os"
	"math/rand"
	"strconv"
	"time"
)

type Point struct{ x, y float64 }

func (p Point) X() float64{
	return p.x
}

func (p Point) Y() float64{
	return p.y
}

// traditional function
func Distance(p, q Point) float64 {
	return math.Hypot(q.X()-p.X(), q.Y()-p.Y())
}

// same thing, but as a method of the Point type
func (p Point) Distance(q Point) float64 {
	return math.Hypot(q.X()-p.X(), q.Y()-p.Y())
}

//!-point

//!+path

// A Path is a journey connecting the points with straight lines.
type Path []Point

// Distance returns the distance traveled along the path.
func (path Path) Distance() float64 {
	sum := 0.0
	fmt.Print("   - ")
	for i := range path {
		if i > 0 {
			fmt.Print(path[i-1].Distance(path[i]), " + ")
			sum += path[i-1].Distance(path[i])
		}
	}
	fmt.Print(path[len(path) - 1].Distance(path[0]))
	sum += path[len(path) - 1].Distance(path[0])
	fmt.Print(" = ")
	return sum
}

func onSegment(p, q, r Point) bool {
	if q.X() <= math.Max(q.X(), r.X()) && q.X() >= math.Min(q.X(), r.X()) &&
		q.Y() <= math.Max(q.Y(), r.Y()) && q.Y() >= math.Min(q.Y(), r.Y()) {
			return true;
	}
	return false;
}

func orientation(p, q, r Point) int{
    val := ( (q.Y() - p.Y()) * (r.X() - p.X()) - (q.X() - p.X()) * (r.Y() - q.Y()) )
	if val == 0 {
		return 0
	}
	if val > 0 {
		return 1
	} else {
		return 2
	}
} 
 
func intersect(p1, q1, p2, q2 Point) bool{ 

    o1 := orientation(p1, q1, p2)
    o2 := orientation(p1, q1, q2) 
    o3 := orientation(p2, q2, p1) 
    o4 := orientation(p2, q2, q1) 
  

    if o1 != o2 && o3 != o4 {	
        return true
    }
  

    if o1 == 0 && onSegment(p1, p2, q1){
    	return true 
    }   
  
    if o2 == 0 && onSegment(p1, q2, q1){
    	return true 
    }
  
    if o3 == 0 && onSegment(p2, p1, q2){
    	return true 
    } 
  
    if o4 == 0 && onSegment(p2, q1, q2){
    	return true 
    }
  
    return false
}
func genVertices (Paths Path, sides int) []Point {
	for i := 0; i < sides; i++ {
    Paths[i].x = (genRange(-100, 100))
    Paths[i].y = (genRange(-100, 100))
    fmt.Println("   - ( ", Paths[i].X(), ", ", Paths[i].Y(), ")")
  }
	return Paths
}

func genRange(min, max int ) float64 {
	seed := rand.NewSource(time.Now().UnixNano())
	rand := rand.New(seed)
	return rand.Float64()*float64((max - min) + min)
}

func main(){
	if len(os.Args) < 2 {
		fmt.Printf("How to compile: ./geometry <number>\n")
	}
	numSides, error := strconv.Atoi(os.Args[1])
	if error == nil {
    	fmt.Printf("- Generating a [%v] sides figure\n", numSides)
	}
	Paths := make(Path, numSides)
	Paths = genVertices(Paths, numSides)
	for intersect(Paths[0], Paths[1], Paths[2], Paths[3]) {
		Paths = genVertices(Paths, numSides)
	}
	fmt.Printf("- Figure's Perimeter\n")
	fmt.Println(Paths.Distance())

}

