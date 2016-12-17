package dirgraph_test

import "testing"
import "fmt"
import "dirgraph"


type testvalues struct {
  input string
  expvalues []int
}

var ourtests = []testvalues{
  { "Graph: AB5, BC4, CD8, DC8, DE6, AD5, CE2, EB3, AE7",
    []int{9, 5, 13, 22, -1, 2, 3, 9, 9, 7}},
  { "", []int{-1, -1, -1, -1, -1, 0, 0, -1, -1, 0}},
  { "Graph: AD1 DC2 CA3", []int{-1, 1, 3, -1, -1, 1, 0, 3, -1, 4}},
  { "AB5 AC3 AD7 BA7 BC5 CB4 DC8", []int{ 10, 7, 15, -1, -1, 2, 5, 3, 9, 10}},
}

var ourdistvalues = [][] string {
    {"A", "B", "C"},
    {"A", "D"},
    {"A", "D", "C"},
    {"A", "E", "B", "C", "D"},
    {"A", "E", "D"},
}

func TDirGraphErr(t *testing.T, testdata string, testno int, expectvalue int,
    receivedvalue int) {
    if expectvalue != receivedvalue {
        t.Error(fmt.Sprintf(
            "For testdata %s and test no %d: expected %d; got %d.\n",
            testdata, testno, expectvalue, receivedvalue))
    }
}

func TestDirGraph(t *testing.T) {
    for _, item := range ourtests {
        ourDG := dirgraph.NewDirectedGraph(item.input)

// We do the route distances.

        for i := 0; i < 4; i++ {
            ourroutelength := ourDG.GetRouteDistance(ourdistvalues[i])
            ourexpected := item.expvalues[i]
            TDirGraphErr(t, item.input, i, ourexpected, ourroutelength)
        }
        getnotrips := ourDG.GetNoTrips("C", "C", 3, false)
        gettripsexp :=  item.expvalues[5]
        TDirGraphErr(t, item.input, 5, gettripsexp, getnotrips)
        getnotrips = ourDG.GetNoTrips("A", "C", 4, true)
        gettripsexp =  item.expvalues[6]
        TDirGraphErr(t, item.input, 6, gettripsexp, getnotrips)
        getshortpath := ourDG.GetShortestPath("A", "C")
        getshortpathexp :=  item.expvalues[7]
        TDirGraphErr(t, item.input, 7, getshortpathexp, getshortpath)
        getshortpath = ourDG.GetShortestPath("B", "B")
        getshortpathexp =  item.expvalues[8]
        TDirGraphErr(t, item.input, 8, getshortpathexp, getshortpath)
        gettripsmax := ourDG.GetNoTripsToMaxDist("C", "C", 30)
        gettripsmaxexp :=  item.expvalues[9]
        TDirGraphErr(t, item.input, 9, gettripsmaxexp, gettripsmax)
    }
}
