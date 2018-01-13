use std::fs::File;
use std::io::Read;

#[derive(Debug)]
struct Header {
    start: String,
    steps: i64
}

struct State {
}

struct Program {
    header: Header,
    states: Vec<State>
}

fn parse_header(header: &str) -> Header {
    let mut lines = header.lines();

    let first_line = lines.next().unwrap().to_string();
    let second_line = lines.next().unwrap().to_string();

    let numbers_in_second_line : String = second_line.chars().filter(|c| c.is_numeric()).collect();

    Header {
        start: first_line.chars().skip(first_line.len() - 2).take(1).collect(),
        steps: numbers_in_second_line.parse().unwrap()
    }
}

#[test]
fn parse_header_test() {
    let header = "Begin in state A.
Perform a diagnostic checksum after 12794428 steps.";

    assert_eq!(parse_header(header).start, "A");
    assert_eq!(parse_header(header).steps, 12794428);
}

fn parse_state(state: &str) -> State {
    State {}
}

fn parse(filename : &str) -> Program {
    let mut file = File::open(filename).expect("Can't open file");
    let mut content = String::new();

    file.read_to_string(&mut content).expect("Can't read file");

    let mut paragraphs = content.split("\n\n");

    let header = parse_header(paragraphs.next().unwrap());
    let mut states = vec![];

    for state in paragraphs {
        states.push(parse_state(state));
    }

    Program {
        header: header,
        states: states
    }
}

fn turing(filename : &str) -> i64 {
    let program = parse(filename);

    0
}

fn main() {
    println!("{}", turing("input.txt"));
}
