use parser::Program;
use std::collections::HashSet;

pub struct Tree {
    pub root : Program,
    pub children: Vec<Box<Tree>>,
}

impl Tree {
    pub fn find_root_name(programs : &Vec<Program>) -> String {
        let mut set : HashSet<String> = HashSet::new();

        for p in programs.iter() {
            set.insert(p.name.clone());
        }

        for p in programs.iter() {
            for sp in p.subprograms.iter() {
                set.remove(sp);
            }
        }

        set.iter().nth(0).unwrap().clone()
    }

    pub fn construct(mut programs : &mut Vec<Program>, root_name : &String) -> Self {
        let index   = programs.iter().position(|p| p.name == *root_name ).unwrap();
        let program = programs.remove(index);

        let mut children : Vec<Box<Tree>> = vec![];

        for s in program.subprograms.iter() {
            children.push(Box::new(Tree::construct(&mut programs, s)));
        }

        Tree { root : program, children : children }
    }

    pub fn display(&self, depth : i32) {
        for _ in 0..depth {
            print!(" ");
        }

        println!("{} ({})", self.root.name, self.root.weight);

        for c in self.children.iter() {
            c.display(depth + 2);
        }
    }

    pub fn total_weight(&self) -> i32 {
        let child_weight : i32 = self.children.iter().map(|c| c.total_weight() ).sum();

        self.root.weight + child_weight
    }

    pub fn diss(&self) -> i32 {
        if self.children.len() == 0 { return 0; }

        let c_diss = self.children.iter().map(|c| c.diss()).filter(|v| *v != 0).nth(0);

        println!("{}: {:?}", self.root.name, c_diss);
        match c_diss {
            None => {
                let min : i32 = self.children.iter().map(|c| c.total_weight() ).min().unwrap();
                let max : i32 = self.children.iter().map(|c| c.total_weight() ).max().unwrap();

                if min == max {
                    0
                } else {
                    if self.children.iter().filter(|c| c.total_weight() == min ).count() > self.children.iter().filter(|c| c.total_weight() == max ).count() {
                        self.children.iter().filter(|c| c.total_weight() == max ).nth(0).unwrap().root.weight - (max - min)
                    } else {
                        self.children.iter().filter(|c| c.total_weight() == min ).nth(0).unwrap().root.weight - (min - max)
                    }
                }
            },

            Some(diss) => diss
        }
    }
}

#[test]
fn tree_test() {
    let mut programs : Vec<Program> = vec![
        Program::new("pbga".to_string(), 66, vec![]),
        Program::new("xhth".to_string(), 57, vec![]),
        Program::new("ebii".to_string(), 61, vec![]),
        Program::new("havc".to_string(), 66, vec![]),
        Program::new("ktlj".to_string(), 57, vec![]),
        Program::new("fwft".to_string(), 72, vec!["ktlj".to_string(), "cntj".to_string(), "xhth".to_string()]),
        Program::new("qoyq".to_string(), 66, vec![]),
        Program::new("padx".to_string(), 45, vec!["pbga".to_string(), "havc".to_string(), "qoyq".to_string()]),
        Program::new("tknk".to_string(), 41, vec!["ugml".to_string(), "padx".to_string(), "fwft".to_string()]),
        Program::new("jptl".to_string(), 61, vec![]),
        Program::new("ugml".to_string(), 68, vec!["gyxo".to_string(), "ebii".to_string(), "jptl".to_string()]),
        Program::new("gyxo".to_string(), 61, vec![]),
        Program::new("cntj".to_string(), 57, vec![]),
    ];

    let root_name = Tree::find_root_name(&mut programs);
    let tree = Tree::construct(&mut programs, &root_name);

    tree.display(0);

    let w = tree.diss();

    println!("{}", w);

    assert_eq!(root_name, "tknk");
}
