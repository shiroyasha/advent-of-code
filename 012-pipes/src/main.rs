use std::collections::HashMap;
use std::collections::HashSet;
use std::fs::File;
use std::io::Read;

struct Graph {
    nodes: HashMap<i32, HashSet<i32>>
}

impl Graph {
    fn new() -> Self {
        Graph { nodes : HashMap::new() }
    }

    fn connect(&mut self, a : i32, b : i32) {
        self.nodes.entry(a).or_insert(HashSet::new()).insert(b);
        self.nodes.entry(b).or_insert(HashSet::new()).insert(a);
    }

    fn nodes(&self) -> Vec<&i32> {
        self.nodes.keys().collect()
    }
}

fn parse(input : &str) -> Graph {
    let mut result = Graph::new();

    input.lines().for_each(|l| {
        let mut parts = l.split("<->");

        let a : i32 = parts.next().unwrap().trim().parse().unwrap();

        parts.next().unwrap().split(",").map(|n| n.trim().parse().unwrap()).for_each(|b| {
            result.connect(a, b);
        });
    });

    result
}

fn reachable_from(root: i32, graph : &Graph) -> HashSet<i32> {
    let mut visited : HashSet<i32> = HashSet::new();
    let mut edge : Vec<i32> = vec![root];

    while edge.len() > 0 {
        let node = graph.nodes.get(&edge.pop().unwrap()).unwrap();

        node.iter().for_each(|n| {
            if !visited.contains(n) {
                edge.push(*n);
                visited.insert(*n);
            }
        })
    }

    visited
}

fn group_count(graph : &Graph) -> i32 {
    let mut nodes = graph.nodes().clone();
    let mut groups = 0;

    while nodes.len() > 0 {
        let root = nodes.pop().unwrap();

        let reachable = reachable_from(*root, graph);

        nodes = nodes.iter().filter(|n| { !reachable.contains(n) }).map(|n| *n).collect();

        groups += 1;
    }

    groups
}


fn main() {
    let mut content = String::new();
    let mut file    = File::open("input.txt").expect("no such file");

    file.read_to_string(&mut content).expect("can't read from file");

    let g = parse(&content);

    println!("{}", reachable_from(0, &g).len());
    println!("{}", group_count(&g));
}

#[cfg(test)]
mod test {
    use super::*;

    #[test]
    fn reachable_from_root_test() {
        let input = "0 <-> 2
1 <-> 1
2 <-> 0, 3, 4
3 <-> 2, 4
4 <-> 2, 3, 6
5 <-> 6
6 <-> 4, 5
";
        let g = parse(input);

        println!("{:?}", g.nodes());

        assert_eq!(reachable_from(0, &g).len(), 6);

        assert_eq!(group_count(&g), 2);
    }
}
