use std::fs::File;
use std::io::Read;

#[derive(Debug)]
struct Layer {
    position: i32,
    depth: i32,
}

impl Layer {

    fn is_deadly(&self, delay : i32) -> bool {
        let time = self.position + delay;
        let pos = time % ((self.depth - 1) * 2);

        // println!("{}, {}, {}", self.position, self.depth, pos);
        // println!("{}, {}", time, pos);

        pos == 0
    }

    fn penatly(&self) -> i32 {
        // println!("penalty: {}", self.position * self.depth);

        self.position * self.depth
    }

}

struct Layers {
    layers : Vec<Layer>
}

impl Layers {
    fn new(input: &str) -> Self {
        let mut layers = Vec::new();

        input.lines().for_each(|l| {
            let mut parts = l.split(":").map(|part| part.trim().parse().unwrap());

            let position = parts.next().unwrap();
            let depth    = parts.next().unwrap();

            layers.push(Layer { position : position, depth : depth });
        });

        Layers { layers : layers }
    }

    fn score(&self) -> i32 {
       self.layers
           .iter()
           .filter(|l| l.is_deadly(0))
           .map(|l| l.penatly())
           .sum()
    }

    fn safe_delay(&self) -> i32 {
        let mut delay = 0;

        loop {
            let has_dedly = self.layers.iter().any(|l| l.is_deadly(delay));

            println!("{}", delay);

            if !has_dedly {
                break;
            }

            delay += 1;
        }

        delay
    }
}

#[test]
fn calculate_test() {
    let input = "0: 3
1: 2
4: 4
6: 4
";

    let layers = Layers::new(input);

    assert_eq!(layers.score(), 24);
}

#[test]
fn safe_delay() {
    let input = "0: 3
1: 2
4: 4
6: 4
";

    let layers = Layers::new(input);

    assert_eq!(layers.safe_delay(), 10);
}

fn main() {
    let mut content = String::new();
    let mut file = File::open("input.txt").expect("can't open file");

    file.read_to_string(&mut content).expect("can't read from file");

   println!("{}", Layers::new(&content).score());
   println!("{}", Layers::new(&content).safe_delay());
}
