use std::fs::File;
use std::io::Read;
use std::fmt;

type Pattern = Vec<char>;

struct Rule {
    input: Vec<Pattern>,
    output: Vec<Pattern>,
}

impl fmt::Display for Rule {
    fn fmt(&self, f: &mut fmt::Formatter) -> fmt::Result {
        write!(f, "{:?} => {:?}", self.input.join(&'/'), self.output.join(&'/'))
    }
}

fn parse(filename : &str) -> Vec<Rule> {
    let mut file = File::open(filename).expect("can't open file");
    let mut content = String::new();

    file.read_to_string(&mut content).expect("can't read from file");

    content.lines().map(|line| {
        let mut parts = line.split("=>");

        let input  = parts.next().unwrap().trim().split("/").map(|p| p.chars().collect()).collect();
        let output = parts.next().unwrap().trim().split("/").map(|p| p.chars().collect()).collect();

        Rule { input, output }
    }).collect()
}

fn main() {
    let rules = parse("input.txt");

    for r in rules {
        println!("{}", r);
    }
}
