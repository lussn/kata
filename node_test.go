package diameter

import (
	"bufio"
	"os"
	"strings"
	"testing"
)

type edge struct{ a, b nodeName }
type edgeList []edge

func (e edgeList) build(g *Graph) {
	for _, edge := range e {
		g.addEdge(edge.a, edge.b)
	}
}

func TestDiameter(t *testing.T) {

	tests := []struct {
		name        string
		edgeList    edgeList
		expDiameter int
	}{
		{
			name: "empty",
		},
		{
			name:        "1 edge",
			edgeList:    edgeList{{"a", "b"}},
			expDiameter: 1,
		},
		{
			name:        "3 in line",
			edgeList:    edgeList{{"a", "b"}, {"b", "c"}},
			expDiameter: 2,
		},
		{
			name:        "4 in line",
			edgeList:    edgeList{{"a", "b"}, {"b", "c"}, {"c", "d"}},
			expDiameter: 3,
		},
		{
			name:        "Triangle",
			edgeList:    edgeList{{"a", "b"}, {"b", "c"}, {"a", "c"}},
			expDiameter: 1,
		},
		{
			name:        "Square",
			edgeList:    edgeList{{"a", "b"}, {"b", "c"}, {"c", "d"}, {"a", "d"}},
			expDiameter: 2,
		},
		{
			name:        "2 loops",
			edgeList:    edgeList{{"a", "b"}, {"b", "c"}, {"c", "a"}, {"c", "d"}, {"d", "e"}, {"e", "c"}},
			expDiameter: 2,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			g := New()
			test.edgeList.build(g)
			dia := g.diameter()
			if dia != test.expDiameter {
				t.Errorf("Diameter not as expected. Have %d, expected %d", dia, test.expDiameter)
			}
		})
	}
}

func BenchmarkDiameter(b *testing.B) {
	g := New()
	// Load the test data
	f, err := os.Open("testdata/edges.txt")
	if err != nil {
		b.Errorf("Could not open file: %s", err)
		return
	}
	defer f.Close()
	s := bufio.NewScanner(f)
	for s.Scan() {
		line := s.Text()
		edge := strings.Fields(line)
		if len(edge) == 0 { // Skip empty lines
			continue
		}
		if len(edge) != 2 {
			b.Error("Expected length of edges to be two")
			return
		}
		g.addEdge(nodeName(edge[0]), nodeName(edge[1]))
	}

	b.Run("diameter", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			d := g.diameter()
			if d != 8000 {
				b.Errorf("Expected diameter to be %d was %d", 8000, d)
			}
		}
	})
}
