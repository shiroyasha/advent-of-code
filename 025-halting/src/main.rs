use std::fs::File;
use std::io::Read;

#[derive(Debug)]
struct Header {
    start: String,
    steps: i64
}

#[derive(Debug, PartialEq, Eq)]
enum Direction {
    Left,
    Right
}

struct Action {
    write: i64,
    movement: Direction,
    continue_with: String
}

struct State {
    name: String,
    action_on_zero: Action,
    action_on_one: Action
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

fn parse_action(lines: &mut std::str::Lines) -> Action {
    lines.next(); // skip first line

    let first_line = lines.next().unwrap();
    let write : String = first_line.chars().skip(first_line.len()-2).take(1).collect();

    lines.next(); // TODO

    let last_line = lines.next().unwrap();
    let continue_with = last_line.chars().skip(first_line.len()-2).take(1).collect();

    Action {
        write: write.parse().unwrap(),
        movement: Direction::Left,
        continue_with: continue_with
    }
}

fn parse_state(state: &str) -> State {
    let mut lines = state.lines();

    let first_line = lines.next().unwrap();
    let name = first_line.chars().skip(first_line.len()-2).take(1).collect();

    let action0 = parse_action(&mut lines);
    let action1 = parse_action(&mut lines);

    State {
        name: name,
        action_on_zero: action0,
        action_on_one: action1
    }
}

#[test]
fn parse_state_test() {
    let state = "In state A:
  If the current value is 0:
    - Write the value 1.
    - Move one slot to the right.
    - Continue with state B.
  If the current value is 1:
    - Write the value 0.
    - Move one slot to the left.
    - Continue with state F.";

    assert_eq!(parse_state(state).name, "A");

    assert_eq!(parse_state(state).action_on_zero.write, 1);
    assert_eq!(parse_state(state).action_on_zero.movement, Direction::Left);
    assert_eq!(parse_state(state).action_on_zero.continue_with, "B");

    assert_eq!(parse_state(state).action_on_one.write, 0);
    assert_eq!(parse_state(state).action_on_one.movement, Direction::Left);
    assert_eq!(parse_state(state).action_on_one.continue_with, "F");
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
