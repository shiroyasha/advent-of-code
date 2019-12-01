use std::io::Read;
use std::fs::File;

#[derive(PartialEq, Debug)]
enum Direction {
    Up,
    Left,
    Right,
    Down
}

fn load_map(filename : &str) -> Vec<Vec<char>> {
    let mut file    = File::open(filename).expect("Can't open file");
    let mut content = String::new();

    let mut result = Vec::new();

    file.read_to_string(&mut content).expect("Can't read from file");

    content.lines().for_each(|line| {
        result.push(line.chars().collect());
    });

    result
}

fn route(filename : &str) -> (String, i64) {
    let map : Vec<Vec<char>> = load_map(filename);

    let mut path = String::new();

    let mut x = map[0].iter().position(|c| *c == '|').unwrap();
    let mut y = 0;
    let mut dir = Direction::Down;
    let mut steps = 0;

    while map[y][x] != ' ' {
        steps += 1;

        println!("{},{}", x, y);

        if map[y][x].is_alphabetic() {
            path.push(map[y][x]);
        }

        match dir {
            Direction::Up => {
                if map[y-1][x] != ' ' {
                    y -= 1;
                } else if map[y][x-1] != ' ' {
                    x -= 1;
                    dir = Direction::Left;
                } else if map[y][x+1] != ' ' {
                    x += 1;
                    dir = Direction::Right;
                } else {
                    break;
                }
            },

            Direction::Down => {
                if map[y+1][x] != ' ' {
                    y += 1;
                } else if map[y][x-1] != ' ' {
                    x -= 1;
                    dir = Direction::Left;
                } else if map[y][x+1] != ' ' {
                    x += 1;
                    dir = Direction::Right;
                } else {
                    break;
                }
            },

            Direction::Right => {
                if map[y][x+1] != ' ' {
                    x += 1;
                } else if map[y-1][x] != ' ' {
                    y -= 1;
                    dir = Direction::Up;
                } else if map[y+1][x] != ' ' {
                    y += 1;
                    dir = Direction::Down;
                } else {
                    break;
                }
            },

            Direction::Left => {
                if map[y][x-1] != ' ' {
                    x -= 1;
                } else if map[y-1][x] != ' ' {
                    y -= 1;
                    dir = Direction::Up;
                } else if map[y+1][x] != ' ' {
                    y += 1;
                    dir = Direction::Down;
                } else {
                    break;
                }
            }
        }
    }

    (path, steps)
}

fn main() {
    println!("{:?}", route("input.txt"));
}

#[test]
fn route_test() {
    assert_eq!(route("test_input.txt"), ("ABCDEF".to_string(), 38));
}
