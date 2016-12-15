#!/usr/bin/env python
# -*- coding: utf-8 -*-
# The trains program (Python implementation).
# Copyright (c) Peter Murphy 2016
# Executed as:
#
# python trains.py [graphfile.txt]
#
# Where graphfile.txt contains a single line representing a directed graph.
# (The syntax of graphfile.txt is explained below.)
# Without command line arguments, the program runs in testing mode.
#
# Rules as explained:
# The local commuter railroad services a number of towns in Kiwiland.  Because
# of monetary concerns, all of the tracks are 'one-way’. That is, a route from
# Kaitaia to Invercargill does not imply the existence of a route from
# Invercargill to Kaitaia. In fact, even if both of these routes do happen
# to exist, they are distinct and are not necessarily the same distance!
#
# The purpose of this problem is to help the railroad provide its customers
# with information about the routes. In particular, you will compute the
# distance along a certain route, the number of different routes between
# two towns, and the shortest route between two towns.
#
# Input: A directed graph where a node represents a town and an edge represents
# a route between two towns. The weighting of the edge represents the distance
# between the two towns. A given route will never appear more than once, and for
# a given route, the starting and ending town will not be the same town.
#
# Output: For test input 1 through 5, if no such route exists, output 'NO SUCH
# ROUTE'. Otherwise, follow the route as given; do not make any extra stops! For
# example, the first problem means to start at city A, then travel directly to
# city B (a distance of 5), then directly to city C (a distance of 4).
#
# 1. The distance of the route A-B-C.
# 2. The distance of the route A-D.
# 3. The distance of the route A-D-C.
# 4. The distance of the route A-E-B-C-D.
# 5. The distance of the route A-E-D.
# 6. The number of trips starting at C and ending at C with a maximum of 3
#   stops.  In the sample data below, there are two such trips: C-D-C (2 stops).
#   and C-E-B-C (3 stops).
# 7. The number of trips starting at A and ending at C with exactly 4 stops.
#   In the sample data below, there are three such trips: A to C (via B,C,D); A
#   to C (via D,C,D); and A to C (via D,E,B).
# 8. The length of the shortest route (in terms of distance to travel) from A
#   to C.
# 9. The length of the shortest route (in terms of distance to travel) from B
#   to B.
# 10. The number of different routes from C to C with a distance of less than
#   30. In the sample data, the trips are: CDC, CEBC, CEBCDC, CDCEBC, CDEBC,
#   CEBCEBC, CEBCEBCEBC.
#
# Test Input:
# For the test input, the towns are named using the first few letters of the
# alphabet from A to D. A route between two towns (A to B) with a distance of 5
# is represented as AB5.
#
# Graph: AB5, BC4, CD8, DC8, DE6, AD5, CE2, EB3, AE7
#
# Expected Output:
# Output #1: 9
# Output #2: 5
# Output #3: 13
# Output #4: 22
# Output #5: NO SUCH ROUTE
# Output #6: 2
# Output #7: 3
# Output #8: 9
# Output #9: 9
# Output #10: 7

import unittest
import sys
import io
import operator

class DirectedGraph:
    """ The DirectedGraph class represents directed graphs: nodes with edges
    between them going one way. Each edge has a weight. For any two nodes, it is
    possible to have two edges between them - going in opposite directions -
    with different weights.

    As for how the graph is specified, see below...
    """

    def __init__(self, strgraphspec):
        """ This method takes a string (strgraphspec), and specifies the whole
        graph from it. Basically, strgraphspec is a string consisting of
        several pieces with the following structure:

        X1Y1n1, X2Y2n2, ...

        Where X1, X2, ... and Y1, Y2, ... are single characters representing
        nodes; and n1, n2, ... are the weight of the edges between them. For
        example:

        AQ1, QA2, AR3

        Represents a graph with three nodes (A, Q and R), one edge from A to Q
        with a weight of 1, one edge from Q to A with a weight of 2, and one
        edge from A to R with a weight of 3.
        """

        self.contents = {}

# The self.contents member is a dictionary/associative array, where keys are
# single letters representing nodes, and values are associative arrays that
# represent the adjacency list for the nodes. For the second sort of arrays,
# keys are the desination nodes adjacent to the first node, and values are the
# weigths of the edges going towards them. This is all we need.
#
# To explain this, let us use our example "AQ1, QA2, AR3". The self.contents
# associative array has two keys "A" and "Q". For "A", the value consists of
# another associate array with "Q" as a key with 1 as a value; and "R" as a key
# with 3 as a value. For the key "Q" in self.contents, the matching value is
# another associative array with "A" as the key and 2 as the value. In other
# words, self.contents is {'A': {'Q': 1, 'R': 3}, 'Q': {'A': 2}}

        graphbits = strgraphspec.replace(',',' ').replace('Graph:', '').split()

# For the safe side, if there is any error in parsing the input string, we make
# self.contents empty.

        try:
            for graphitem in graphbits:
                nodebegin = graphitem[0]
                nodeend = graphitem[1]
                edgeweight = int(graphitem[2:])
                if nodebegin not in self.contents:
                    self.contents[nodebegin] = {nodeend: edgeweight}
                else:
                    self.contents[nodebegin][nodeend] = edgeweight
        except:
            self.contents = {}

    def getRouteDistance(self, nodearray):
        """ This gets the total distance travelling along nodearray[0],
        nodearray[1], and on to nodearray[len - 1]. (Here, nodearray is a list
        of characters representing nodes.) If no such path exists, the method
        returns None.
        """
        if len(nodearray) == 0 or nodearray[0] not in self.contents:
            return None
        firstpairnode = nodearray[0]
        sumsofar = 0
        for otherpairnode in nodearray[1:]:
            if otherpairnode not in self.contents[firstpairnode]:
                return None
            sumsofar +=  self.contents[firstpairnode][otherpairnode]
            firstpairnode = otherpairnode
        return sumsofar


    def getNoTrips(self, startnode, endnode, nostops, isequals):
        """ This counts the number of trips possible from one node to another
        up to (or equal to) a number of stops. The parameters:
        startnode: the node that is the start of the trip (a character)
        endnode: the node that is the end of the trip (a character)
        nostops: the maximum number of stops involved
        isequals: if True, then the only trips counted are when the number of
        stops is equal to nostops. If False, then trips are counted when the
        number of stops is less than or equal to nostops.
        """
        countstops = 0
        if startnode not in self.contents:
            return countstops
        for othernode in self.contents[startnode]:
            if (othernode == endnode) and ((isequals and (nostops == 1))
                    or (not isequals)):
                countstops += 1
            if nostops > 1:
                countstops += self.getNoTrips(othernode, endnode,
                        nostops - 1, isequals)
        return countstops

    def getNoTripsToMaxDist(self, startnode, endnode, distancebound):
        """ This counts the number of trips possible from one node to another
        below a distance bound. The parameters:
        startnode: the node that is the start of the trip (a character)
        endnode: the node that is the end of the trip (a character)
        distancebound: trips with distances underneath it are counted; trips
        with distances equal to or exceeding this are not counted
        """
        countstops = 0
        if distancebound < 0:
            return countstops
        if startnode not in self.contents:
            return countstops
        for othernode in self.contents[startnode]:
            edgelength = self.contents[startnode][othernode]
            if (othernode == endnode) and edgelength < distancebound:
                countstops += 1
            countstops += self.getNoTripsToMaxDist(othernode, endnode,
                    distancebound - edgelength)
        return countstops

    def getShortestPath(self, startnode, endnode):
        """ This calculates the distance of the shortest path from startnode to
        endnode. This is based on Dijkstra's algorithm. If a shortest path is
        found, the distance is returned. When no path exists, None is returned.
        """
        if startnode not in self.contents or endnode not in self.contents:
            return None
        distmap = {} # A map of distances.
        infinity = float('inf')

# We start by initializing the distances to those adjacent to startnode.

        for node in self.contents:
            if node in self.contents[startnode]:
                distmap[node] = self.contents[startnode][node]
            else:
                distmap[node] = infinity # Represents infinity

# Rather than have a set of visited nodes and removing them, we basically have
# a sequence (oursort_distmap) of nodes sorted by distance, and an index,
# indexofmintocalc, which refers to them.

        indexofmintocalc = 0
        while len(distmap) != indexofmintocalc:
            oursort_distmap = sorted(distmap.items(),
                key=operator.itemgetter(1))
            minnode = oursort_distmap[indexofmintocalc][0]
            minnodelength = oursort_distmap[indexofmintocalc][1]

# If we come to a situation where the distance in the sorted nodes is infinity,
# then we can say that this node (and other nodes) are unapproachable. We can
# break the loop with good conscience.

            if minnodelength == infinity:
                break;
            indexofmintocalc += 1
            for neighbournode in self.contents[minnode]:
                if ((self.contents[minnode][neighbournode] + minnodelength) \
                    < distmap[neighbournode]):
                    distmap[neighbournode] = \
                            self.contents[minnode][neighbournode] \
                            + minnodelength
        if distmap[endnode] == infinity:
            return None
        return distmap[endnode]


class TestGraphs(unittest.TestCase):

    def test_examplegraph(self):
        ourDG = DirectedGraph("AB5, BC4, CD8, DC8, DE6, AD5, CE2, EB3, AE7")
        self.assertEqual(ourDG.getRouteDistance(["A", "B", "C"]), 9)
        self.assertEqual(ourDG.getRouteDistance(["A", "D"]), 5)
        self.assertEqual(ourDG.getRouteDistance(["A", "D", "C"]), 13)
        self.assertEqual(ourDG.getRouteDistance(["A", "E", "B", "C", "D"]), 22)
        self.assertEqual(ourDG.getRouteDistance(["A", "E", "D"]), None)
        self.assertEqual(ourDG.getNoTrips("C", "C", 3, False), 2)
        self.assertEqual(ourDG.getNoTrips("A", "C", 4, True), 3)
        self.assertEqual(ourDG.getShortestPath("A", "C"), 9)
        self.assertEqual(ourDG.getShortestPath("B", "B"), 9)
        self.assertEqual(ourDG.getNoTripsToMaxDist("C", "C", 30), 7)

    def test_emptygraph(self):
        ourDG = DirectedGraph("")
        self.assertEqual(ourDG.getRouteDistance(["A", "B", "C"]), None)
        self.assertEqual(ourDG.getRouteDistance(["A", "D"]), None)
        self.assertEqual(ourDG.getRouteDistance(["A", "D", "C"]), None)
        self.assertEqual(ourDG.getRouteDistance(["A", "E", "B", "C", "D"]),
            None)
        self.assertEqual(ourDG.getRouteDistance(["A", "E", "D"]), None)
        self.assertEqual(ourDG.getNoTrips("C", "C", 3, False), 0)
        self.assertEqual(ourDG.getNoTrips("A", "C", 4, True), 0)
        self.assertEqual(ourDG.getShortestPath("A", "C"), None)
        self.assertEqual(ourDG.getShortestPath("B", "B"), None)
        self.assertEqual(ourDG.getNoTripsToMaxDist("C", "C", 30), 0)

    def test_simplegraph(self):
        ourDG = DirectedGraph("Graph: AD1 DC2 CA3")
        self.assertEqual(ourDG.getRouteDistance(["A", "B", "C"]), None)
        self.assertEqual(ourDG.getRouteDistance(["A", "D"]), 1)
        self.assertEqual(ourDG.getRouteDistance(["A", "D", "C"]), 3)
        self.assertEqual(ourDG.getRouteDistance(["A", "E", "B", "C", "D"]),
            None)
        self.assertEqual(ourDG.getRouteDistance(["A", "E", "D"]), None)
        self.assertEqual(ourDG.getNoTrips("C", "C", 3, False), 1)
        self.assertEqual(ourDG.getNoTrips("A", "C", 4, True), 0)
        self.assertEqual(ourDG.getNoTrips("A", "C", 4, False), 1)
        self.assertEqual(ourDG.getShortestPath("A", "C"), 3)
        self.assertEqual(ourDG.getShortestPath("B", "B"), None)
        self.assertEqual(ourDG.getNoTripsToMaxDist("C", "C", 30), 4)
        self.assertEqual(ourDG.getNoTripsToMaxDist("C", "C", 31), 5)

    def test_oddgraph(self):
        ourDG = DirectedGraph("AB5 AC3 AD7 BA7 BC5 CB4 DC8")
        self.assertEqual(ourDG.getRouteDistance(["A", "B", "C"]), 10)
        self.assertEqual(ourDG.getRouteDistance(["A", "D"]), 7)
        self.assertEqual(ourDG.getRouteDistance(["A", "D", "C"]), 15)
        self.assertEqual(ourDG.getRouteDistance(["A", "E", "B", "C", "D"]),
            None)
        self.assertEqual(ourDG.getRouteDistance(["A", "E", "D"]), None)
        self.assertEqual(ourDG.getNoTrips("C", "C", 3, False), 2)
        self.assertEqual(ourDG.getNoTrips("A", "C", 4, True), 5)
        self.assertEqual(ourDG.getShortestPath("A", "C"), 3)
        self.assertEqual(ourDG.getShortestPath("B", "B"), 9)
        self.assertEqual(ourDG.getNoTripsToMaxDist("C", "C", 30), 10)


if __name__ == '__main__':

# With no command line arguments, assume that the program is in testing mode.
# With one argument, assume that this is a file with a single line containing
# the graph information.

    def prettyprint(index, inputin):
        if inputin is not None:
            print("Output #{index}: {inputin}".format(**vars()))
        else:
            print("Output #{index}: NO SUCH ROUTE".format(**vars()))

    if len(sys.argv) == 1:
        unittest.main()
    else:
        with io.open(sys.argv[1], "r", encoding="utf-8") as f:
            ourline = f.readline()
            ourDirGraph = DirectedGraph(ourline)
            prettyprint(1, ourDirGraph.getRouteDistance(["A", "B", "C"]))
            prettyprint(2, ourDirGraph.getRouteDistance(["A", "D"]))
            prettyprint(3, ourDirGraph.getRouteDistance(["A", "D", "C"]))
            prettyprint(4, ourDirGraph.getRouteDistance(["A", "E", "B", "C",
                "D"]))
            prettyprint(5, ourDirGraph.getRouteDistance(["A", "E", "D"]))
            prettyprint(6, ourDirGraph.getNoTrips("C", "C", 3, False))
            prettyprint(7, ourDirGraph.getNoTrips("A", "C", 4, True))
            prettyprint(8, ourDirGraph.getShortestPath("A", "C"))
            prettyprint(9, ourDirGraph.getShortestPath("B", "B"))
            prettyprint(10, ourDirGraph.getNoTripsToMaxDist("C", "C", 30))
