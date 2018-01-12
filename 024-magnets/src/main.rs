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

fn dfs(start : i64, magnets: Vec<&Magnet>) -> (i64, i64) {
    let mut connectable = vec![];

    for m in magnets.iter() {
       if m.a == start || m.b == start {
           connectable.push(m.clone());
       }
    }

    let size = connectable.len();

    if size == 0 {
        return (0, 0);
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
              let (len, strength) = dfs(m.b, others);

              (len+1, strength + m.value())
          } else {
              // print!("--{}/{}", m.a, m.b);
              let (len, strength) = dfs(m.a, others);

              (len+1, strength + m.value())
          }
       }).max_by(|a, b| {
           if a.0 == b.0 {
               a.1.cmp(&b.1)
           } else {
               a.0.cmp(&b.0)
           }
       }).unwrap()
    }
}

fn calculate(filename : &str) -> i64 {
    let magnets = parse(filename);

    for m in magnets.iter() {
        println!("{}/{}", m.a, m.b);
    }

    let (_, strength) = dfs(0, magnets.iter().collect());

    strength
}

#[test]
fn calculate_test() {
    assert_eq!(calculate("test_input.txt"), 19);
}

fn main() {
    println!("{}", calculate("input.txt"));
}
