use std::fs::File;
use std::io::Read;

struct Magnet {
    index: usize,
    a: i64,
    b: i64
}

impl Magnet {
    fn value(&self) -> i64 {
        self.a + self.b
    }
}

fn parse(filename : &str) -> Vec<Magnet> {
    let mut file = File::open(filename).expect("Can't open file");
    let mut content = String::new();

    file.read_to_string(&mut content).expect("can't read from file");

    content.lines().enumerate().map(|(index, line)| {
        let mut parts = line.split("/");

        let a = parts.next().unwrap().parse().unwrap();
        let b = parts.next().unwrap().parse().unwrap();

        Magnet { index: index, a: a, b: b }
    }).collect()
}

fn dfs(start : i64, magnets: Vec<&Magnet>) -> i64 {
    let mut connectable = vec![];

    for m in magnets.iter() {
       if m.a == start || m.b == start {
           connectable.push(m.clone());
       }
    }

    let size = connectable.len();

    if size == 0 {
        return 0;
    } else {
       connectable.iter().map(|m| {
           let mut others = vec![];

           for m2 in magnets.iter() {
               if m.index != m2.index {
                   others.push(m2.clone());
               }
           }

          if m.a == start {
              // print!("--{}/{}", m.a, m.b);
              m.value() + dfs(m.b, others)
          } else {
              // print!("--{}/{}", m.a, m.b);
              m.value() + dfs(m.a, others)
          }
       }).max().unwrap()
    }
}

fn calculate(filename : &str) -> i64 {
    let magnets = parse(filename);

    for m in magnets.iter() {
        println!("{}/{}", m.a, m.b);
    }

    dfs(0, magnets.iter().collect())
}

#[test]
fn calculate_test() {
    assert_eq!(calculate("test_input.txt"), 31);
}

fn main() {
    println!("{}", calculate("input.txt"));
}
