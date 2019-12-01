// use std::fs::File;
// use std::io::Read;
// use std::collections::HashMap;

// type Commands = Vec<Command>;

// enum Command {
//     Set(String, String),
//     Sub(String, String),
//     Mul(String, String),
//     Jnz(String, String)
// }

// fn parse(filename : &str) -> Commands {
//     let mut file = File::open(filename).expect("can't open file");
//     let mut content = String::new();

//     file.read_to_string(&mut content).expect("can't read from file");

//     content.lines().map(|line| {
//         let mut parts = line.split_whitespace();

//         let cmd = parts.next().unwrap();
//         let a = parts.next().unwrap();
//         let b = parts.next().unwrap();

//         match cmd {
//             "set" => Command::Set(a.to_string(), b.to_string()),
//             "mul" => Command::Mul(a.to_string(), b.to_string()),
//             "sub" => Command::Sub(a.to_string(), b.to_string()),
//             "jnz" => Command::Jnz(a.to_string(), b.to_string()),
//             _ => panic!("unrecognized patter")
//         }
//     }).collect()
// }

// fn execute(filename : &str) -> i64 {
//     let commands = parse(filename);
//     let mut current : i64 = 0;
//     let mut registers : HashMap<String, i64> = HashMap::new();

//     registers.insert("a".to_string(), 1);

//     loop {
//         for (a, b) in registers.iter() {
//             print!("{}: {}, ", a, b);
//         }
//         println!("");

//         match commands.get(current as usize).unwrap() {
//             &Command::Set(ref a, ref b) => {
//                 println!("set {} {}", a, b);

//                 let val_b = match b.parse() {
//                     Ok(number) => {
//                         number
//                     },
//                     Err(_) => {
//                         match registers.get(b) {
//                             Some(v) => *v,
//                             None => 0
//                         }
//                     }
//                 };

//                 registers.insert(a.clone(), val_b);

//                 current += 1;
//             },

//             &Command::Mul(ref a, ref b) => {
//                 println!("mul {} {}", a, b);

//                 let val_a = match a.parse() {
//                     Ok(number) => {
//                         number
//                     },
//                     Err(_) => {
//                         match registers.get(a) {
//                             Some(v) => *v,
//                             None => 0
//                         }
//                     }
//                 };

//                 let val_b = match b.parse() {
//                     Ok(number) => {
//                         number
//                     },
//                     Err(_) => {
//                         match registers.get(b) {
//                             Some(v) => *v,
//                             None => 0
//                         }
//                     }
//                 };

//                 registers.insert(a.clone(), val_a * val_b);

//                 current += 1;
//             },

//             &Command::Sub(ref a, ref b) => {
//                 println!("sub {} {}", a, b);

//                 let val_a = match a.parse() {
//                     Ok(number) => {
//                         number
//                     },
//                     Err(_) => {
//                         match registers.get(a) {
//                             Some(v) => *v,
//                             None => 0
//                         }
//                     }
//                 };

//                 let val_b = match b.parse() {
//                     Ok(number) => {
//                         number
//                     },
//                     Err(_) => {
//                         match registers.get(b) {
//                             Some(v) => *v,
//                             None => 0
//                         }
//                     }
//                 };

//                 registers.insert(a.clone(), val_a - val_b);

//                 current += 1;
//             },

//             &Command::Jnz(ref a, ref b) => {
//                 println!("jnz {} {}", a, b);

//                 let val_a = match a.parse() {
//                     Ok(number) => {
//                         number
//                     },
//                     Err(_) => {
//                         match registers.get(a) {
//                             Some(v) => *v,
//                             None => 0
//                         }
//                     }
//                 };

//                 let val_b = match b.parse() {
//                     Ok(number) => {
//                         number
//                     },
//                     Err(_) => {
//                         match registers.get(b) {
//                             Some(v) => *v,
//                             None => 0
//                         }
//                     }
//                 };

//                 if val_a != 0 {
//                     current += val_b;
//                 } else {
//                     current += 1;
//                 }
//             }
//         }

//         if current < 0 || current >= (commands.len() as i64) {
//             break;
//         }
//     }

//     match registers.get("h") {
//         Some(v) => *v,
//         None => 0
//     }
// }

fn calculate() -> i64 {
    let mut h = 0;

    let mut b = 79 * 100 + 100_000;
    let limit = b + 17_000;

    loop {
      for d in 2..b {
          if b % d == 0 {
              h += 1;
              break;
          }
      }

      b += 17;
      if b > limit { break; }
    }

    h
}

fn main() {
    // println!("{}", execute("input.txt"));
    println!("{}", calculate());
}
