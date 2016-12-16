/**
  * @desc The trains program (JavaScript implementation).
  * Executed as:
  *
  * nodejs trains.js [graphfile.txt]
  *
  * Where graphfile.txt contains a single line representing a directed graph.
  * (The syntax of graphfile.txt is explained in trains.txt.)
  * Without command line arguments, the program runs in testing mode.
  *
  * @author Peter Murphy peterkmurphy@gmail.com
  * @required unit.js
*/

"use strict";
var fs = require ("fs");
var test = require('unit.js');

/**
  * @desc sorts an object by its property values (in ascending order)
  * @param object $objectin - the object for property value sorting
  * @return Array of items, where item[0] is property name and item[1] is
  * property value.
*/

function sortobjbyprops (objectin) {
    var arrayforsorting = [];
    for (var item in objectin) {
        arrayforsorting.push([item, objectin[item]]);
    }
    arrayforsorting.sort(function(theone, theother)
    {
        return theone[1]-theother[1];
    });
    return arrayforsorting;
}

/**
  * @desc The DirectedGraph class represents directed graphs: nodes with edges
  * between them going one way. Each edge has a weight. For any two nodes, it is
  * possible to have two edges between them - going in opposite directions -
  * with different weights. As for how the graph is specified, see the
  * constructor
*/

class DirectedGraph {

/**
  * @desc Constructs a DirectedGraph.
  * @param string $strgraphspec. The string consists of several pieces with the
  * following structure:
  *
  * X1Y1n1, X2Y2n2, ...
  *
  * Where X1, X2, ... and Y1, Y2, ... are single characters representing
  * nodes; and n1, n2, ... are the weight of the edges between them. For
  * example:
  *
  * AQ1, QA2, AR3
  *
  * Represents a graph with three nodes (A, Q and R), one edge from A to Q
  * with a weight of 1, one edge from Q to A with a weight of 2, and one
  * edge from A to R with a weight of 3.
*/

    constructor(strgraphspec) {

/**
  * The this.contents member is a dictionary/associative array, where keys are
  *  single letters representing nodes, and values are associative arrays that
  * represent the adjacency list for the nodes. For the second sort of arrays,
  * keys are the desination nodes adjacent to the first node, and values are the
  * weigths of the edges going towards them. This is all we need.
  *
  * To explain this, let us use our example "AQ1, QA2, AR3". The this.contents
  * associative array has two keys "A" and "Q". For "A", the value consists of
  * another associate array with "Q" as a key with 1 as a value; and "R" as a
  * key with 3 as a value. For the key "Q" in this.contents, the matching value
  * is another associative array with "A" as the key and 2 as the value. In
  * other words, this.contents is {'A': {'Q': 1, 'R': 3}, 'Q': {'A': 2}}.
*/

        this.contents = {};
        var graphbits = strgraphspec.replace('Graph:', '').split(/[ ,]+/);
        if (graphbits[0] == '') {
            graphbits.shift(); // If the first item is an empty string, remove.
        }
        // If there are errors parsing the input string, make contents empty
        try {
            for (var i = 0; i < graphbits.length; i++) {
                var graphitem = graphbits[i];
                var nodebegin = graphitem[0];
                var nodeend = graphitem[1];
                var edgeweight = parseInt(graphitem.slice(2), 10);
                if (!(nodebegin in this.contents)) {
                    var added = {};
                    added[nodeend] = edgeweight;
                    this.contents[nodebegin] = added;
                } else {
                    this.contents[nodebegin][nodeend] = edgeweight;
                }
            }
        } catch (e) {
            this.contents = {};
        }
    }

    /**
      * @desc gets the total distance travelling along a route of nodes.
      * @param array[string] $nodearray - a list of nodes, where each value is
      * a character representing a node.
      * @return int (the distance) or null if no route exists.
    */

    getRouteDistance(nodearray) {
        if (nodearray.length == 0) {
            return null;
        }
        var firstpairnode = nodearray[0];
        if (!(firstpairnode in this.contents)) {
            return null;
        }
        var sumsofar = 0;
        for (var i = 1; i < nodearray.length; i++) {
            var otherpairnode = nodearray[i];
            if (!(otherpairnode in this.contents[firstpairnode])) {
                return null;
            }
            sumsofar +=  this.contents[firstpairnode][otherpairnode];
            firstpairnode = otherpairnode;
        }
        return sumsofar;
    }

    /**
      * @desc his counts the number of trips possible from one node to another
      * up to (or equal to) a number of stops.
      * @param string $startnode - the starting node (a character)
      * @param string $endnode - the ending node (a character)
      * @param int $nostops - the maximum number of stops involved
      * @param bool $isequals - if true, then the only trips counted are when
      * the number of stops is equal to nostops. If false, then trips are
      * counted when the number of stops is less than or equal to nostops.
      * @return int - the number of routes
    */

    getNoTrips(startnode, endnode, nostops, isequals) {
        var countstops = 0;
        if (!(startnode in this.contents)) {
            return countstops;
        }
        for (var othernode in this.contents[startnode]) {
            if ((othernode == endnode) && ((isequals && (nostops == 1))
                    || (!isequals))) {
                countstops += 1;
            }
            if (nostops > 1) {
                countstops += this.getNoTrips(othernode, endnode,
                        nostops - 1, isequals);
            }
        }
        return countstops;
    }

    /**
      * @desc Counts the number of trips possible from one node to another
      * below a distance bound
      * @param string $startnode - the starting node (a character).
      * @param string $endnode - the ending node (a character).
      * @param int $distancebound - the distance bound. Trips with distance
      * underneath are counted; trips with distances equal to or exceeding are
      * not counted.
      * @return int - the number of trips
    */

    getNoTripsToMaxDist(startnode, endnode, distancebound) {
        var countstops = 0
        if (distancebound <= 0) {
            return countstops;
        }
        if (!(startnode in this.contents)) {
            return countstops;
        }
        for (var othernode in this.contents[startnode]) {
            var edgelength = this.contents[startnode][othernode];
            if ((othernode == endnode) && (edgelength < distancebound)) {
                countstops += 1;
            }
            countstops += this.getNoTripsToMaxDist(othernode, endnode,
                      distancebound - edgelength);
        }
        return countstops;
    }

    /**
      * @desc finds the distance of the shortest path between two nodes. This is
      * based on Dijkstra's algorithm.
      * @param string $startnode - the starting node (a character).
      * @param string $endnode - the ending node (a character).
      * @return int (the distance of the shortest path) or null if no path
      * exists.
    */

    getShortestPath(startnode, endnode) {
        if (!(startnode in this.contents)){
            return null;
        }
        if (!(endnode in this.contents)){
            return null;
        }
        var distmap = {};

// We start by initializing the distances to those adjacent to startnode.

        for (var node in this.contents) {
            if (node in this.contents[startnode]) {
                distmap[node] = this.contents[startnode][node];
            }
            else {
                distmap[node] = Infinity;
            }
        }
        var distmaplength = Object.keys(distmap).length;
/*
Rather than have a set of visited nodes and removing them, we basically have
a sequence (oursort_distmap) of nodes sorted by distance, and an index,
indexofmintocalc, which refers to them.
*/

        var indexofmintocalc = 0;
        while (distmaplength != indexofmintocalc) {
            var oursort_distmap = sortobjbyprops(distmap);
            var minnode = oursort_distmap[indexofmintocalc][0];
            var minnodelength = oursort_distmap[indexofmintocalc][1];

/*
If we come to a situation where the distance in the sorted nodes is infinity,
then we can say that this node (and other nodes) are unapproachable. We can
break the loop with good conscience.
*/
            if (minnodelength == Infinity) {
                break;
            }
            indexofmintocalc += 1
            for (var neighbournode in this.contents[minnode]) {
                if ((this.contents[minnode][neighbournode] + minnodelength)
                    < distmap[neighbournode]) {
                    distmap[neighbournode] =
                            this.contents[minnode][neighbournode]
                            + minnodelength;
                }
            }
        }
        if (distmap[endnode] == Infinity) {
            return null;
        }
        return distmap[endnode];
    }
}

/**
  * @desc pretty prints the output required for the trains program
  * @param int $index - the line number
  * @param $inputin - the input. May be a string, an int or null
  * @return doesn't return anything. Just prints strings of the form:
  * "Output #index: inputin"
  * or if inputin is null:
  * "Output #index: NO SUCH ROUTE"
*/

function prettyprint(count, inputin) {
    if (!(inputin == null)) {
        console.log(`Output #${count}: ${inputin}`);
    }
    else {
        console.log(`Output #${count}: NO SUCH ROUTE`);
    }
}

/**
  * @desc A quick bit of code to get the first line of a file
  * @param string $fileinput - the file name
  * @return the first line (or an empty string if the file is empty).
*/

function getfirstline(fileinput) {
    var fileoutput = fs.readFileSync(fileinput, 'utf8');
    fileoutput = fileoutput.split('\n')[0] // Just get the first line
    return fileoutput;
}

/**
  * This is the "main" part of the program. Now command line arguments are
  * looked at. If a filename is passed it, it is read and processed. Otherwise
  * the program goes into testing mode.
*/

if (process.argv.length > 2) {
    var fileinput = process.argv[2];
    var ourDirGraph = new DirectedGraph(getfirstline(fileinput));
    prettyprint(1, ourDirGraph.getRouteDistance(["A", "B", "C"]));
    prettyprint(2, ourDirGraph.getRouteDistance(["A", "D"]));
    prettyprint(3, ourDirGraph.getRouteDistance(["A", "D", "C"]));
    prettyprint(4, ourDirGraph.getRouteDistance(["A", "E", "B", "C",
        "D"]));
    prettyprint(5, ourDirGraph.getRouteDistance(["A", "E", "D"]));
    prettyprint(6, ourDirGraph.getNoTrips("C", "C", 3, false));
    prettyprint(7, ourDirGraph.getNoTrips("A", "C", 4, true));
    prettyprint(8, ourDirGraph.getShortestPath("A", "C"));
    prettyprint(9, ourDirGraph.getShortestPath("B", "B"));
    prettyprint(10, ourDirGraph.getNoTripsToMaxDist("C", "C", 30));
}
else {
    console.log("Here are some tests!");

    /**
      * @desc A function to test the DirectedGraph class.
      * @param array[string] $va - an array of values. The first is the file
      * name is the test data. The next 10 should be expected values for the
      * method calls required by trains.txt in order. (The first is the
      * routine distance A-B-C, and so on to the number of trips from C to C
      * under a distance of 30.)
      * @return there is no return value. The code has an assertion if
      * the expect values don't match the obtained values.
    */


    function testDirGraph(va) {
        console.log(`We are testing ${va[0]}.`);
        var ourDG = new DirectedGraph(getfirstline(va[0]));
        test.value(ourDG.getRouteDistance(["A", "B", "C"])).isEqualTo(va[1]);
        test.value(ourDG.getRouteDistance(["A", "D"])).isEqualTo(va[2]);
        test.value(ourDG.getRouteDistance(["A", "D", "C"])).isEqualTo(va[3]);
        test.value(ourDG.getRouteDistance(["A", "E", "B", "C", "D"]))
            .isEqualTo(va[4]);
        test.value(ourDG.getRouteDistance(["A", "E", "D"])).isEqualTo(va[5]);
        test.value(ourDG.getNoTrips("C", "C", 3, false)).isEqualTo(va[6]);
        test.value(ourDG.getNoTrips("A", "C", 4, true)).isEqualTo(va[7]);
        test.value(ourDG.getShortestPath("A", "C")).isEqualTo(va[8]);
        test.value(ourDG.getShortestPath("B", "B")).isEqualTo(va[9]);
        test.value(ourDG.getNoTripsToMaxDist("C", "C", 30))
            .isEqualTo(va[10]);
    }

    var maindata = ['graph1.txt', 9, 5, 13, 22, null, 2, 3, 9, 9, 7];
    testDirGraph(maindata);
    var zerodata = ['graph2.txt', null, null, null, null, null, 0, 0, null,
        null, 0];
    testDirGraph(zerodata);
    var cyclicdata = ['graph3.txt', null, 1, 3, null, null, 1, 0, 3, null, 4];
    testDirGraph(cyclicdata);
    var odddata = ['graph4.txt', 10, 7, 15, null, null, 2, 5, 3, 9, 10];
    testDirGraph(cyclicdata);
    console.log("All the tests have passed!");
}
