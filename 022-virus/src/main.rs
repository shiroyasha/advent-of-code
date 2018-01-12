use std::fs::File;
use std::io::Read;

#[derive(Clone)]
enum Direction {
    Up,
    Down,
    Left,
    Right
}

#[derive(Clone)]
enum State {
    Clean,
    Weakened,
    Infected,
    Flagged
}

struct Map {
    field: Vec<Vec<State>>,
    x: usize,
    y: usize,
    direction: Direction,
}

impl Map {
    fn turn_right(&mut self) {
        match self.direction {
            Direction::Up    => self.direction = Direction::Right,
            Direction::Down  => self.direction = Direction::Left,
            Direction::Left  => self.direction = Direction::Up,
            Direction::Right => self.direction = Direction::Down,
        }
    }

    fn turn_left(&mut self) {
        match self.direction {
            Direction::Up    => self.direction = Direction::Left,
            Direction::Down  => self.direction = Direction::Right,
            Direction::Left  => self.direction = Direction::Down,
            Direction::Right => self.direction = Direction::Up,
        }
    }

    fn advance(&mut self) {
        match self.direction {
            Direction::Up    => self.y -= 1,
            Direction::Down  => self.y += 1,
            Direction::Left  => self.x -= 1,
            Direction::Right => self.x += 1,
        }
    }
}

fn parse(filename : &str) -> Map {
    let mut file = File::open(filename).expect("Can't open file");
    let mut content = String::new();

    file.read_to_string(&mut content).expect("Can't read from file");

    let mut field : Vec<Vec<State>> = Vec::with_capacity(1_000);

    println!("A");

    for _ in 0..1_000 {
        let mut row : Vec<State> = vec![];
        row.resize(1_000, State::Clean);

        field.push(row);
    }

    println!("B");

    content.lines().enumerate().for_each(|(i, line)| {
        line.chars().enumerate().for_each(|(j, chr)| {
            if chr == '#' {
                field[500 + i][500 + j] = State::Infected;
            } else {
                field[500 + i][500 + j] = State::Clean;
            }
        });
    });

    println!("C");

    let height : usize = content.lines().count();
    let width  : usize = content.lines().next().unwrap().chars().count();

    println!("D");

    Map {
        field: field,
        x: 500 + width/2,
        y: 500 + height/2,
        direction: Direction::Up
    }
}

fn walk(filename : &str) -> i64 {
    let mut map = parse(filename);
    let mut infection_count = 0;

    for index in 0..10_000_000 {
        // for i in 490..510 {
        //     for j in 490..510 {
        //         let symbol = match map.field[i][j] {
        //             State::Clean => '.',
        //             State::Weakened => 'W',
        //             State::Infected => '#',
        //             State::Flagged => 'F',
        //         };

        //         if i == map.y && j == map.x {
        //             print!("[{}]", symbol);
        //         } else {
        //             print!(" {} ", symbol);
        //         }
        //     }
        //     println!("");
        // }

        // println!("-------");

        if index % 10_000 == 0 {
            println!("{}, {}, {}", index, map.x, map.y);
        }

        match map.field[map.y][map.x] {
          State::Clean => {
              map.turn_left();
              map.field[map.y][map.x] = State::Weakened;
          },
          State::Weakened => {
              map.field[map.y][map.x] = State::Infected;
              infection_count += 1;
          },
          State::Infected => {
              map.turn_right();
              map.field[map.y][map.x] = State::Flagged;
          },
          State::Flagged => {
              map.turn_right();
              map.turn_right();

              map.field[map.y][map.x] = State::Clean;
          }
        }

        map.advance()
    }

    infection_count
}

#[test]
fn walk_test() {
    let infection_count = walk("test_input.txt");

    assert_eq!(infection_count, 2511944);
}

fn main() {
    println!("{}", walk("input.txt"));
}
