use std::fs::File;
use std::io::prelude::*;

mod parser;
mod tree;

fn main() {
    let mut file = File::open("input.txt").expect("Can't open file");
    let mut content = String::new();

    file.read_to_string(&mut content).expect("Can't read file");

    let mut programs = parser::Program::parse_all(&content);

    let root_name = tree::Tree::find_root_name(&programs);
    let tree = tree::Tree::construct(&mut programs, &root_name);

    tree.display(9);

    let dissbalance  = tree.diss();

    println!("part 1: {}", tree.root.name);
    println!("part 2: {}", dissbalance);
}
