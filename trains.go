//usr/bin/go run $0 $@ ; exit
// The trains program (Go implementation).
// Copyright (c) Peter Murphy 2016
// Executed as:
//
// go run trains.go [graphfile.txt]
//
// Where graphfile.txt contains a single line representing a directed graph.
// (The syntax of graphfile.txt is explained below.)
// Without command line arguments, the program does nothing.
//
// Rules as explained:
// The local commuter railroad services a number of towns in Kiwiland.  Because
// of monetary concerns, all of the tracks are 'one-way’. That is, a route from
// Kaitaia to Invercargill does not imply the existence of a route from
// Invercargill to Kaitaia. In fact, even if both of these routes do happen
// to exist, they are distinct and are not necessarily the same distance!
//
// The purpose of this problem is to help the railroad provide its customers
// with information about the routes. In particular, you will compute the
// distance along a certain route, the number of different routes between
// two towns, and the shortest route between two towns.
//
// Input: A directed graph where a node represents a town and an edge represents
// a route between two towns. The weighting of the edge represents the distance
// between the two towns. A given route will never appear more than once, and
// for a given route, the starting and ending town will not be the same town.
//
// Output: For test input 1 through 5, if no such route exists, output 'NO SUCH
// ROUTE'. Otherwise, follow the route as given; do not make any extra stops!
// For example, the first problem means to start at city A, then travel directly
// to city B (a distance of 5), then directly to city C (a distance of 4).
//
// 1. The distance of the route A-B-C.
// 2. The distance of the route A-D.
// 3. The distance of the route A-D-C.
// 4. The distance of the route A-E-B-C-D.
// 5. The distance of the route A-E-D.
// 6. The number of trips starting at C and ending at C with a maximum of 3
//   stops.  In the sample data below, there are two such trips: C-D-C (2 stops)
//   and C-E-B-C (3 stops).
// 7. The number of trips starting at A and ending at C with exactly 4 stops.
//   In the sample data below, there are three such trips: A to C (via B,C,D); A
//   to C (via D,C,D); and A to C (via D,E,B).
// 8. The length of the shortest route (in terms of distance to travel) from A
//   to C.
// 9. The length of the shortest route (in terms of distance to travel) from B
//   to B.
// 10. The number of different routes from C to C with a distance of less than
//   30. In the sample data, the trips are: CDC, CEBC, CEBCDC, CDCEBC, CDEBC,
//   CEBCEBC, CEBCEBCEBC.
//
// Test Input:
// For the test input, the towns are named using the first few letters of the
// alphabet from A to D. A route between two towns (A to B) with a distance of 5
// is represented as AB5.
//
// Graph: AB5, BC4, CD8, DC8, DE6, AD5, CE2, EB3, AE7
//
// Expected Output:
// Output #1: 9
// Output #2: 5
// Output #3: 13
// Output #4: 22
// Output #5: NO SUCH ROUTE
// Output #6: 2
// Output #7: 3
// Output #8: 9
// Output #9: 9
// Output #10: 7

package main
import (
    "fmt"
    "io/ioutil"
    "os"
    "sort"
    "strings"
    "strconv"
    "unicode/utf8"
)

// The following code is adapted from:
// http://stackoverflow.com/questions/18695346/how-to-sort-a-mapstringint-by-its-values
// It is renamed here.

// The SortMapByValues sorts a map of strings to integers by its values (in
// ascending order). It returns an array of Pairs, where Key is for the original
// keys, and Value is for the original values.
func SortMapByValues(mapstringtoint map[string]int) PairList{
  pl := make(PairList, len(mapstringtoint))
  i := 0
  for k, v := range mapstringtoint {
    pl[i] = Pair{k, v}
    i++
  }
  sort.Sort(pl)
  return pl
}

type Pair struct {
  Key string
  Value int
}

type PairList []Pair

func (p PairList) Len() int { return len(p) }
func (p PairList) Less(i, j int) bool { return p[i].Value < p[j].Value }
func (p PairList) Swap(i, j int){ p[i], p[j] = p[j], p[i] }

// This the end of the adapted code. Now the original code begins.

// The DirectedGraph class represents directed graphs: nodes with edges between
// them going one way. Each edge has a weight. For any two nodes, it is possible
// to have two edges between them - going in opposite directions - with
// different weights. The NewDirectedGraph function represents the constructor.
//
// The contents member is a mapping, where keys are single letters representing
// nodes, and values are other maps that represent the adjacency list for the
// nodes. For the second sort of maps, keys are the desination nodes adjacent to
//  the first node, and values are the weigths of the edges going towards them.
// This is all we need.
//
// To explain this, let us use an example "AQ1, QA2, AR3". The self.contents
// associative array has two keys "A" and "Q". For "A", the value consists of
// another associate array with "Q" as a key with 1 as a value; and "R" as a key
// with 3 as a value. For the key "Q" in self.contents, the matching value is
// another associative array with "A" as the key and 2 as the value. In other
// words, self.contents is {'A': {'Q': 1, 'R': 3}, 'Q': {'A': 2}}
type DirectedGraph struct {
    contents  map[string] map[string] int
}

type DirGraphEdge struct {
    sourcenode string
    destnode string
    edgelength int
}

func ParseEdge(stredgespec string) *DirGraphEdge {
    if utf8.RuneCount([]byte(stredgespec)) < 3 {
        return nil;
    }
    dge := DirGraphEdge{}
    runes_array := []rune(stredgespec)
    dge.sourcenode = string(runes_array[0])
    dge.destnode = string(runes_array[1])
    i, err := strconv.Atoi(string(runes_array[2:]))
    if err != nil {
        return nil;
    }
    dge.edgelength = i
    return &dge
}

// The NewDirectedGraph function takes a string (strgraphspec), and specifies
// a whole DirectedGraph object from it. Basically, strgraphspec is a string
// consisting of several pieces with the following structure:
//
// X1Y1n1, X2Y2n2, ...
//
// Where X1, X2, ... and Y1, Y2, ... are single characters representing nodes;
// and n1, n2, ... are the weight of the edges between them. For example:
//
// AQ1, QA2, AR3
// Represents a graph with three nodes (A, Q and R), one edge from A to Q with a
// weight of 1, one edge from Q to A with a weight of 2, and one edge from A to
// R with a weight of 3.
func NewDirectedGraph(strgraphspec string) *DirectedGraph {
    dg := DirectedGraph{}
    dg.contents = make(map[string] map[string] int)
    strreplace := strings.Replace(strgraphspec, ","," ", -1)
    strreplace = strings.Replace(strreplace, "Graph:", "", -1)
    strarray := strings.Fields(strreplace)
    for i := 0; i < len(strarray); i++ {
        stredge := ParseEdge(strarray[i])
        if stredge != nil {
            elem, ok := dg.contents[stredge.sourcenode]
            if ok == true {
                elem[stredge.destnode] = stredge.edgelength
            } else {
                ourmap := make(map[string]int)
                ourmap[stredge.destnode] = stredge.edgelength
                dg.contents[stredge.sourcenode] = ourmap
            }
        }
    }
    return &dg
}

// The getRouteDistance method gets the total distance travelling along
// nodearray[0], nodearray[1], and on to nodearray[len - 1]. (Here, nodearray is
// a list of characters representing nodes.) If no such path exists, the method
// returns -1.
func (dg DirectedGraph) getRouteDistance(nodearray []string) int {
    if len(nodearray) == 0 {
        return -1
    }
    firstpairnode, ok := dg.contents[nodearray[0]]
    if !ok {
        return -1
    }
    sumsofar := 0
    for _, otherpairnode := range nodearray[1:] {
        otherpairlength, otherok := firstpairnode[otherpairnode]
        if !otherok {
            return -1
        }
        sumsofar +=  otherpairlength
        firstpairnode = dg.contents[otherpairnode]
    }
    return sumsofar
}

// The getNoTrips method This counts the number of trips possible from one node
// to another up to (or equal to) a number of stops. The parameters:
// startnode: the node that is the start of the trip (a character)
// endnode: the node that is the end of the trip (a character)
// nostops: the maximum number of stops involved
// isequals: if true, then the only trips counted are when the number of stops
// is equal to nostops. If false, then trips are counted when the
// number of stops is less than or equal to nostops.
func (dg DirectedGraph) getNoTrips(startnode string, endnode string,
        nostops int, isequals bool) int {
    countstops := 0
    startnodecont, ok := dg.contents[startnode]
    if ok == false {
        return countstops
    }
    for othernode := range startnodecont {
        if (othernode == endnode) && ((isequals && (nostops == 1)) ||
                (! isequals)) {
            countstops += 1
        }
        if nostops > 1 {
            countstops += dg.getNoTrips(othernode, endnode,
                    nostops - 1, isequals)
        }
    }
    return countstops
}

func (dg DirectedGraph) getNoTripsToMaxDist(startnode string, endnode string,
        distancebound int) int {
    countstops := 0
    if (distancebound < 0) {
        return countstops
    }
    startnodecont, ok := dg.contents[startnode]
    if ok == false {
        return countstops
    }
    for othernode := range startnodecont {
        edgelength := startnodecont[othernode]
        if (othernode == endnode) && edgelength < distancebound {
            countstops += 1
        }
        countstops += dg.getNoTripsToMaxDist(othernode, endnode,
                distancebound - edgelength)
    }
    return countstops
}

func (dg DirectedGraph) getShortestPath(startnode string, endnode string) int {
    startnodecont, ok := dg.contents[startnode]
    if ok == false {
        return -1
    }
    _, otherok :=  dg.contents[endnode]
    if otherok == false {
        return -1
    }
    distmap := make(map[string]int) // A map of distances.
    const MaxUint = ^uint(0)
    const infinity = int(MaxUint >> 1)

// We start by initializing the distances to those adjacent to startnode.

    for node := range dg.contents {
        immediateedge, thisok := startnodecont[node]
        if thisok {
            distmap[node] = immediateedge
        } else {
            distmap[node] = infinity
        }
    }

// Rather than have a set of visited nodes and removing them, we basically have
// a sequence (oursort_distmap) of nodes sorted by distance, and an index,
// indexofmintocalc, which refers to them.

    indexofmintocalc := 0
    for len(distmap) != indexofmintocalc {
        oursort_distmap := SortMapByValues(distmap)
        minnode := oursort_distmap[indexofmintocalc].Key
        minnodelength := oursort_distmap[indexofmintocalc].Value

// If we come to a situation where the distance in the sorted nodes is infinity,
// then we can say that this node (and other nodes) are unapproachable. We can
// break the loop with good conscience.

        if minnodelength == infinity {
            break
        }
        indexofmintocalc += 1
        for neighbournode := range dg.contents[minnode] {
            if ((dg.contents[minnode][neighbournode] + minnodelength) <
                    distmap[neighbournode]) {
                distmap[neighbournode] = dg.contents[minnode][neighbournode] +
                    minnodelength
            }
        }
    }
    if (distmap[endnode] == infinity) {
        return -1
    }
    return distmap[endnode]
}


func prettyprint(index int, inputina int) {
    if inputina != -1 {
        fmt.Print(fmt.Sprintf("Output #%d: %d\n", index, inputina))
    } else {
        fmt.Print(fmt.Sprintf("Output #%d: NO SUCH ROUTE\n", index))
    }
}


func main() {
    argsWithoutProg := os.Args[1:]
    if (len(argsWithoutProg) == 1) {
        filename := argsWithoutProg[0]
        dat, err := ioutil.ReadFile(filename)
        if err != nil {
            fmt.Print("No such file or directory, I'm afraid.\n")
            os.Exit(0)
        }
        firstline := strings.Split(string(dat), "/n")[0]
        ourDG := NewDirectedGraph(firstline)
        prettyprint(1, ourDG.getRouteDistance([]string{"A", "B", "C"}))
        prettyprint(2, ourDG.getRouteDistance([]string{"A", "D"}))
        prettyprint(3, ourDG.getRouteDistance([]string{"A", "D", "C"}))
        prettyprint(4,
            ourDG.getRouteDistance([]string{"A", "E", "B", "C", "D"}))
        prettyprint(5, ourDG.getRouteDistance([]string{"A", "E", "D"}))
        prettyprint(6, ourDG.getNoTrips("C", "C", 3, false))
        prettyprint(7, ourDG.getNoTrips("A", "C", 4, true))
        prettyprint(8, ourDG.getShortestPath("A", "C"))
        prettyprint(9, ourDG.getShortestPath("B", "B"))
        prettyprint(10, ourDG.getNoTripsToMaxDist("C", "C", 30))
    }
}
