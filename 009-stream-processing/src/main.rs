use std::fs::File;
use std::io::Read;

const STATE_START : i32 = 1;
const STATE_GROUP_START : i32 = 2;
const STATE_GROUP_OVER : i32 = 3;
const STATE_GARBAGE_START : i32 = 4;
const STATE_GARBAGE_OVER : i32 = 5;
const STATE_GARBAGE_IGNORE_NEXT : i32 = 6;
const STATE_COMMA : i32 = 7;

fn debug(state: i32, c : char) {
    let name = match state {
        STATE_START => "START",
        STATE_GROUP_START => "GROUP_START",
        STATE_GROUP_OVER => "GROUP_OVER",
        STATE_GARBAGE_START => "GARBAGE_START",
        STATE_GARBAGE_OVER => "GARBAGE_OVER",
        STATE_GARBAGE_IGNORE_NEXT => "GARBGE_IGNORE_NEXT",
        STATE_COMMA => "COMMA",
        _ => "INVALID"
    };

    println!("char: {}, state: {}", c, name);
}

fn parser(content : &str) -> (i32, i32) {
    let mut state = STATE_START;
    let mut depth = 0;
    let mut sum = 0;
    let mut garbage = 0;

    for c in content.chars() {
        debug(state, c);

        match state {
            STATE_START => match c {
                '{' => {
                    state = STATE_GROUP_START;
                    depth += 1;
                },

                _   => panic!("Invalid character")
            },

            STATE_GROUP_START => match c {
                '{' => {
                    state = STATE_GROUP_START;
                    depth += 1;
                },

                '}' => {
                    state = STATE_GROUP_OVER;
                    sum += depth;
                    depth -= 1;
                },

                '<' => {
                    state = STATE_GARBAGE_START;
                },

                _ => panic!("Invalid character")
            },

            STATE_GARBAGE_START => match c {
                '>' => {
                    state = STATE_GARBAGE_OVER;
                },

                '!' => {
                    state = STATE_GARBAGE_IGNORE_NEXT;
                },

                _ => {
                    state = STATE_GARBAGE_START;
                    garbage += 1;
                }
            },

            STATE_GARBAGE_IGNORE_NEXT => match c {
                _ => {
                    state = STATE_GARBAGE_START;
                }
            },

            STATE_GARBAGE_OVER | STATE_GROUP_OVER => match c {
                ',' => {
                    state = STATE_COMMA;
                },

                '}' => {
                    state = STATE_GROUP_OVER;
                    sum += depth;
                    depth -= 1;
                },

                _ => panic!("Invalid character")
            },

            STATE_COMMA => match c {
                '{' => {
                    state = STATE_GROUP_START;
                    depth += 1;
                },

                '<' => {
                    state = STATE_GARBAGE_START;
                },

                _ => panic!("Invalid character")
            },

            _ => panic!("Invalid state")
        }
    }

    (sum, garbage)
}

#[test]
fn parser_test() {
    assert_eq!(parser("{}"), (1, 0));
    assert_eq!(parser("{{{}}}"), (6, 0));
    assert_eq!(parser("{{},{}}"), (5, 0));
    assert_eq!(parser("{{{},{},{{}}}}"), (16, 0));
    assert_eq!(parser("{<a>,<a>,<a>,<a>}"), (1, 4));
    assert_eq!(parser("{{<ab>},{<ab>},{<ab>},{<ab>}}"), (9, 8));
    assert_eq!(parser("{{<!!>},{<!!>},{<!!>},{<!!>}}"), (9, 0));
    assert_eq!(parser("{{<a!>},{<a!>},{<a!>},{<ab>}}"), (3, 16));
}

fn main() {
    let mut content = String::new();
    let mut file = File::open("input.txt").expect("no file");

    file.read_to_string(&mut content).expect("can't read from file");

    println!("sum: {:?}", parser(content.trim()));
}
