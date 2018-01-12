use std::fs::File;
use std::io::Read;

enum Direction {
    Up,
    Down,
    Left,
    Right
}

struct Map {
    field: Vec<Vec<i8>>,
    x: usize,
    y: usize,
    direction: Direction,
}

impl Map {
    fn is_infected(&self) -> bool {
        self.field[self.y][self.x] == 1
    }

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

    fn clean(&mut self) {
        self.field[self.y][self.x] = 0;
    }

    fn infect(&mut self) {
        self.field[self.y][self.x] = 1;
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

    let mut field = Vec::with_capacity(20_000);

    println!("A");

    for _ in 0..20_000 {
        let mut row = vec![];
        row.resize(20_000, 0);

        field.push(row);
    }

    println!("B");

    content.lines().enumerate().for_each(|(i, line)| {
        line.chars().enumerate().for_each(|(j, chr)| {
            if chr == '#' {
                field[10_000 + i][10_000 + j] = 1;
            } else {
                field[10_000 + i][10_000 + j] = 0;
            }
        });
    });

    println!("C");

    let height : usize = content.lines().count();
    let width  : usize = content.lines().next().unwrap().chars().count();

    println!("D");

    Map {
        field: field,
        x: 10_000 + width/2,
        y: 10_000 + height/2,
        direction: Direction::Up
    }
}

fn walk(filename : &str) -> i64 {
    let mut map = parse(filename);
    let mut infection_count = 0;

    for _ in 0..10_000 {
        // for i in 9990..10020 {
        //     for j in 9990..10020 {
        //         if i == map.y && j == map.x {
        //             if map.field[i][j] == 1 {
        //                 print!("[#]");
        //             } else {
        //                 print!("[.]");
        //             }
        //         } else {
        //             if map.field[i][j] == 1 {
        //                 print!(" # ");
        //             } else {
        //                 print!(" . ");
        //             }
        //         }
        //     }
        //     println!("");
        // }

        println!("{}, {}", map.x, map.y);
        if map.is_infected() {
           map.turn_right();
           map.clean();
        } else {
           map.turn_left();
           map.infect();
            infection_count += 1;
        }

        map.advance()
    }


    infection_count
}

#[test]
fn walk_test() {
    let infection_count = walk("test_input.txt");

    assert_eq!(infection_count, 5587);
}

fn main() {
    println!("{}", walk("input.txt"));
}
