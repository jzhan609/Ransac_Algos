/*
Name: Jacob Zhang
Course: CSI2120
ID: 300231094
Date: 3/11/2023
Description: A RANSAC Algorithm written in GoLang.
*/
package main

// Imports
import (
	"bufio"
	"fmt"
	"math"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

// Point3D struct of variables x, y, z
type Point3D struct {
	X float64
	Y float64
	Z float64
}

// Plane3D struct of variables a, b, c, d
type Plane3D struct {
	A float64
	B float64
	C float64
	D float64
}

// Plane3DwSupport struct to assign a support size to a Plane3D
type Plane3DwSupport struct {
	Plane3D
	SupportSize int
}

// Reads a file using bufio
func ReadXYZ(filename string) ([]Point3D, error) {

	// Checks for errors in fileName
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	var points []Point3D

	// Scans each line and appends it to the list of points
	scanner := bufio.NewScanner(file)
	c := 0

	// Loops while there's more lines
	for scanner.Scan() {

		// Reads the line, splits it, and appends it as a Point3D to the Point3D array
		line := scanner.Text()
		if c != 0 {
			l := strings.Split(line, "\t")
			x, _ := strconv.ParseFloat(l[0], 64)
			y, _ := strconv.ParseFloat(l[1], 64)
			z, _ := strconv.ParseFloat(l[2], 64)

			points = append(points, Point3D{x, y, z})
		}
		c++
	}

	file.Close()

	return points, nil
}

// Saves a slice of Point3D into an XYZ file
func SaveXYZ(filename string, points []Point3D) error {

	// Checks for errors
	file, err := os.Create(filename)
	if err != nil {
		return err
	}

	// Saves to a file
	for _, point := range points {
		_, err := fmt.Fprintf(file, "%.6f\t%.6f\t%.6f\n", point.X, point.Y, point.Z)
		if err != nil {
			return err
		}
	}

	return nil
}

// Calculates distance between 2 points
func (p1 *Point3D) GetDistance(p2 *Point3D) float64 {

	return math.Sqrt(math.Pow((p1.X-p2.X), 2) + math.Pow((p1.Y-p2.Y), 2) + math.Pow((p1.Z-p2.Z), 2))
}

// Calculates the closest distance a Point3D is from a Plane3D
func (plane *Plane3D) GetDistance(p *Point3D) float64 {
	distance := math.Abs(plane.A*p.X+plane.B*p.Y+plane.C*p.Z+plane.D) / math.Sqrt(plane.A*plane.A+plane.B*plane.B+plane.C*plane.C)
	return distance
}

// Gets the plane created from 3 points using cross product
func GetPlane(points [3]Point3D) Plane3D {

	p1 := points[0]
	p2 := points[1]
	p3 := points[2]

	a := (p2.Y-p1.Y)*(p3.Z-p1.Z) - (p3.Y-p1.Y)*(p2.Z-p1.Z)
	b := (p2.Z-p1.Z)*(p3.X-p1.X) - (p3.Z-p1.Z)*(p2.X-p1.X)
	c := (p2.X-p1.X)*(p3.Y-p1.Y) - (p3.X-p1.X)*(p2.Y-p1.Y)
	d := -(a*p1.X + b*p1.Y + c*p1.Z)
	return Plane3D{a, b, c, d}
}

// Gets the number of iterations from confidence and percent values
func GetNumberOfIterations(c float64, p float64) int {
	numIterations := int(math.Log10(1-c) / math.Log10(1-math.Pow(p, 3)))
	return numIterations
}

// Computes the support of a plane in a set of points
func GetSupport(plane Plane3D, points []Point3D, eps float64) Plane3DwSupport {

	supportSize := 0

	// Checks the distance every point is from the plane.
	for _, point := range points {
		if plane.GetDistance(&point) <= eps {
			supportSize++
		}
	}

	return Plane3DwSupport{Plane3D: plane, SupportSize: supportSize}
}

// Returns all of the supporting Point3Ds for a Plane3D
func GetSupportingPoints(plane Plane3D, points []Point3D, eps float64) []Point3D {

	var supPoints []Point3D

	// Checks each Point3D to see if it is a support
	for _, point := range points {
		if plane.GetDistance(&point) <= eps {
			supPoints = append(supPoints, point)
		}
	}

	return supPoints
}

// Removes all the Point3Ds which are part of a Plane3D and returns the remaining Point3Ds
func RemovePlane(plane Plane3D, points []Point3D, eps float64) []Point3D {
	var remainingPoints []Point3D

	for _, point := range points {
		if plane.GetDistance(&point) > eps {
			remainingPoints = append(remainingPoints, point)
		}
	}

	return remainingPoints
}

// Generates a random Point3D from a Point3D array
func randomPointGenerator(points []Point3D) Point3D {

	randIndex := rand.Intn(len(points))
	return points[randIndex]
}

// Repeats a function until told to stop
func repeatFct(wg *sync.WaitGroup, stop <-chan bool, fct func([]Point3D) Point3D, points []Point3D) <-chan Point3D {

	ptStream := make(chan Point3D)

	go func() {

		// Stops the go routine if told to
		defer func() { wg.Done() }()
		defer close(ptStream)

		// Passes the output of a fct to the point stream. Loops until told to stop
		for {
			select {
			case <-stop:
				return
			case ptStream <- fct(points):
			}
		}
	}()

	return ptStream
}

// Generates triplets of Point3Ds
func tripletsGenerator(wg *sync.WaitGroup, stop <-chan bool, PointStream <-chan Point3D) <-chan [3]Point3D {

	tripletStream := make(chan [3]Point3D)

	go func() {

		//Stops the go routine if told to
		var tmpTriplet [3]Point3D
		defer func() { wg.Done() }()
		defer close(tripletStream)

		// Passes triplets to the triplet stream until told to stop
		for {
			select {
			case <-stop:
				return

			default:
				for i := 0; i < 3; i++ {
					tmpTriplet[i] = <-PointStream
				}
				select {
				case <-stop:
					return
				case tripletStream <- tmpTriplet:
				}
			}
		}
	}()

	return tripletStream
}

// Reads arrays of Point3D and resend them. It automatically stops the pipeline after having received N arrays.
func takeN(wg *sync.WaitGroup, stop <-chan bool, inputTripletStream <-chan [3]Point3D, numIterations int) <-chan [3]Point3D {

	outputStream := make(chan [3]Point3D)

	go func() {

		defer func() { wg.Done() }()
		defer close(outputStream)

		// Outputs numIterations number of triplets
		for i := 0; i < numIterations; i++ {
			select {
			case <-stop:
				break
			case outputStream <- <-inputTripletStream:
			}
		}
	}()

	return outputStream
}

// Estimates the plane defined by a triplet
func PlaneEstimator(wg *sync.WaitGroup, stop <-chan bool, inputTripletStream <-chan [3]Point3D) <-chan Plane3D {

	outputStream := make(chan Plane3D)

	go func() {

		defer func() { wg.Done() }()
		defer close(outputStream)

		// For each triplet, outputs the lane defined by the triplet
		for i := range inputTripletStream {
			select {
			case <-stop:
				return
			case outputStream <- GetPlane(i):
			}
		}
	}()

	return outputStream
}

// Takes a Plane3D and Point3D array and outputs a Plane3DwSupport with the plane and number of supports
func SupportingPointsFinder(wg *sync.WaitGroup, stop <-chan bool, inputPlaneStream <-chan Plane3D, points []Point3D, eps float64) <-chan Plane3DwSupport {

	outputStream := make(chan Plane3DwSupport)

	go func() {

		defer func() { wg.Done() }()
		defer close(outputStream)

		// For each plane in the stream, gets the number of supports and outputs it in the output stream
		for i := range inputPlaneStream {
			select {
			case <-stop:
				return
			default:
				support := GetSupport(i, points, eps)
				select {
				case <-stop:
					return
				case outputStream <- support:
				}
			}
		}
	}()

	return outputStream
}

// Multiplexes the results received from multiple channels into one output channel
func fanIn(wg *sync.WaitGroup, stop <-chan bool, channels []<-chan Plane3DwSupport) <-chan Plane3DwSupport {

	var wg2 sync.WaitGroup
	outputStream := make(chan Plane3DwSupport)

	reader := func(ch <-chan Plane3DwSupport) {
		defer func() { wg2.Done() }()
		for i := range ch {

			select {
			case <-stop:
				return
			case outputStream <- i:
			}
		}
	}

	wg2.Add(len(channels))
	for _, ch := range channels {

		go reader(ch)
	}

	go func() {

		defer func() { wg.Done() }()
		defer close(outputStream)
		wg2.Wait()
	}()

	return outputStream
}

// Identifies the dominant plane
func dominantPlaneIdentifier(bestSup *Plane3DwSupport, inputPlaneStream <-chan Plane3DwSupport) {

	// Loops for each Plane3D coming from the plane stream
	for i := range inputPlaneStream {

		// If the plane has a better support than the previous best, it replaces the previous best supported plane
		if i.SupportSize > bestSup.SupportSize {
			*bestSup = i
		}
	}
}

func main() {

	// Identifies the args and sets the randomizer seed
	file := os.Args[1]
	c, _ := strconv.ParseFloat(os.Args[2], 64)
	p, _ := strconv.ParseFloat(os.Args[3], 64)
	eps, _ := strconv.ParseFloat(os.Args[4], 64)
	points, _ := ReadXYZ(file)
	numIterations := GetNumberOfIterations(c, p)
	rand.Seed(time.Now().UnixNano())

	// Loops 3 times
	for j := 1; j <= 3; j++ {

		// Creates the stop channel, wait group, and dummy best support plane
		stop := make(chan bool)
		var wg sync.WaitGroup
		bestSupp := Plane3DwSupport{Plane3D{0, 0, 0, 0}, 0}

		// Runs the pipeline until the PlaneEstimator step
		wg.Add(4)
		randomStream := PlaneEstimator(&wg, stop, takeN(&wg, stop, tripletsGenerator(&wg, stop, repeatFct(&wg, stop, randomPointGenerator, points)), numIterations))
		fanOut := 22

		// Makes 'fanOut' number of threads and finds the supports for each thread
		wg.Add(fanOut)
		filterStreams := make([]<-chan Plane3DwSupport, fanOut)
		for i := 0; i < fanOut; i++ {
			filterStreams[i] = SupportingPointsFinder(&wg, stop, randomStream, points, eps)
		}

		// Identifies the dominant plane
		dominantPlaneIdentifier(&bestSupp, fanIn(&wg, stop, filterStreams))

		// Creates the variables used for generating the output files
		supportingPoints := GetSupportingPoints(bestSupp.Plane3D, points, eps)
		RemovePlane(bestSupp.Plane3D, points, eps)
		outputString := fmt.Sprintf("%s_p%d", file, j)
		outputString2 := fmt.Sprintf("%s_p0%d", file, j)

		// Generates the output files
		SaveXYZ(outputString, supportingPoints)
		SaveXYZ(outputString2, points)
	}
}
