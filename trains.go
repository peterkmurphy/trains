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

/* The following taken from:
http://stackoverflow.com/questions/18695346/how-to-sort-a-mapstringint-by-its-values
*/

func rankByWordCount(wordFrequencies map[string]int) PairList{
  pl := make(PairList, len(wordFrequencies))
  i := 0
  for k, v := range wordFrequencies {
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

/* Now the main code begins */

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
        oursort_distmap := rankByWordCount(distmap)
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

func check(e error) {
    if e != nil {
        fmt.Print("No such file or directory, I'm afraid.")
    }
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
        check(err)
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
