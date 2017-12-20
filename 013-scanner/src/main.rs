use std::collections::HashMap;
use std::fs::File;
use std::io::Read;

#[derive(Debug)]
struct Layer {
    length : i32,
    position: i32,
    forward: bool,
}

impl Layer {
    fn new(length : i32) -> Self {
        Layer { length: length, position : 0, forward : true }
    }

    fn reset(&mut self) {
        self.position = 0;
        self.forward = true;
    }

    fn advance(&mut self) {
        if self.forward {
            self.position += 1;
        } else {
            self.position -= 1;
        }

        if self.position == 0 {
            self.forward = true;
        }

        if self.position == self.length - 1 {
            self.forward = false;
        }
    }
}

fn display(layers : &HashMap<i32, Layer>, current_pos : i32) {
    let layers_len = *layers.keys().max().unwrap() + 1;

    println!("---------------------");

    for i in 0..layers_len {
        if current_pos == i {
            print!("({}): ", i);
        } else {
            print!(" {} : ", i);
        }

        match layers.get(&i) {
            Some(l) => {
                for j in 0..l.length {
                    if l.position == j {
                        print!("[S] ");
                    } else {
                        print!("[ ] ");
                    }
                }

                println!("");
            },
            None => {
                println!("...");
            }
        }
    }
}

fn calculate(start : i32, layers : &mut HashMap<i32, Layer>) -> (i32, bool) {
    let mut score = 0;
    let mut caught = false;
    let mut current_pos : i32 = start;
    let layers_len = *layers.keys().max().unwrap();

    println!(" --> {}", start);

    while current_pos <= layers_len {
        match layers.get(&current_pos) {
            Some(l) => {
                if l.position == 0 {
                    score += current_pos * l.length;
                    caught = true;
                }
            },
            None => {
                ()
            }
        }

        current_pos += 1;

        layers.iter_mut().for_each(|(i, l)| l.advance());
    }

    (score, caught)
}

fn wait_time_to_be_safe(layers : &mut HashMap<i32, Layer>) -> i32 {
    let mut wait = 0;

    loop {
        layers.iter_mut().for_each(|(i, l)| l.reset());

        let caught = calculate(-wait, layers).1;

        if !caught {
            break;
        }

        wait += 1;
    }

    wait
}

#[test]
fn calculate_test() {
   let mut layers : HashMap<i32, Layer> = HashMap::new();

   layers.insert(0, Layer::new(3));
   layers.insert(1, Layer::new(2));
   layers.insert(4, Layer::new(4));
   layers.insert(6, Layer::new(4));

   println!("{:?}", layers);

   assert_eq!(calculate(0, &mut layers).0, 24);
}

#[test]
fn wait_time_to_be_safe_test() {
   let mut layers : HashMap<i32, Layer> = HashMap::new();

   layers.insert(0, Layer::new(3));
   layers.insert(1, Layer::new(2));
   layers.insert(4, Layer::new(4));
   layers.insert(6, Layer::new(4));

   assert_eq!(wait_time_to_be_safe(&mut layers), 10);
}

fn main() {
    let mut content = String::new();
    let mut file = File::open("input.txt").expect("can't open file");

    file.read_to_string(&mut content).expect("can't read from file");

    let mut layers : HashMap<i32, Layer> = HashMap::new();

    content.lines().for_each(|l| {
        let mut parts = l.split(":").map(|part| part.trim().parse().unwrap());

        let index = parts.next().unwrap();
        let length = parts.next().unwrap();

        layers.insert(index, Layer::new(length));
    });


   println!("{}", calculate(0, &mut layers).0);
   // println!("{}",wait_time_to_be_safe(&mut layers));
}
