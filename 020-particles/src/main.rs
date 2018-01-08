extern crate regex;

use std::fs::File;
use std::io::Read;
use regex::Regex;

#[derive(PartialEq, Debug)]
struct Particle {
    pos: (i64, i64, i64),
    vel: (i64, i64, i64),
    acc: (i64, i64, i64),
}

fn parse(filename : &str) -> Vec<Particle> {
    let mut file = File::open(filename).expect("Can't open file");
    let mut content = String::new();

    file.read_to_string(&mut content).expect("Can't read from file");

    let regex = Regex::new(r"p=<(.*),(.*),(.*)>, v=<(.*),(.*),(.*)>, a=<(.*),(.*),(.*)>").unwrap();

    content.lines().map(|line| {
        let captures = regex.captures(line).unwrap();
        let fetch_num = |index| { captures.get(index).unwrap().as_str().trim().parse().unwrap() };

        Particle {
            pos : (fetch_num(1), fetch_num(2), fetch_num(3)),
            vel : (fetch_num(4), fetch_num(5), fetch_num(6)),
            acc : (fetch_num(7), fetch_num(8), fetch_num(9)),
        }
    }).collect()
}

fn main() {
    println!("{:?}", parse("test_input.txt"));
}

#[test]
fn parse_test() {
    assert_eq!(parse("test_input.txt"), vec![
       Particle { pos : (3, 0, 0), vel : (2, 0, 0), acc: (-1, 0, 0) },
       Particle { pos : (4, 0, 0), vel : (0, 0, 0), acc: (-2, 0, 0) },
    ]);
}
