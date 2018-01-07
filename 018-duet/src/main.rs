use std::fs::File;
use std::io::Read;

mod program;

use program::Program;

fn main() {
    let mut content = String::new();
    let mut file = File::open("input.txt").expect("file not found");

    file.read_to_string(&mut content).expect("can't read file");

    let mut p = Program::new(&content);

    p.run();

    println!("{}", p.result);
}

#[test]
fn execute_test() {
    let input = "set a 1
add a 2
mul a a
mod a 5
snd a
set a 0
rcv a
jgz a -1
set a 1
jgz a -2
";

    let mut p = Program::new(&input);

    p.run();

   assert_eq!(p.result, 4);
}
