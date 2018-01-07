use std::collections::HashMap;
use std::fs::File;
use std::io::Read;

fn main() {
    let mut content = String::new();
    let mut file = File::open("input.txt").expect("file not found");

    file.read_to_string(&mut content).expect("can't read file");

    println!("{}", execute(&content));
}

fn execute(input : &str) -> i64 {
    let instructions : Vec<&str> = input.lines().collect();
    let mut registers : HashMap<&str, i64> = HashMap::new();
    let mut sounds : Vec<i64> = Vec::new();
    let mut current : i64 = 0;

    while current >= 0 && current < (instructions.len() as i64) {
        for (k, v) in &registers {
            print!("{}:{} ", k, v);
        }

        println!("\n{}", instructions[current as usize]);

        let mut words = instructions[current as usize].split_whitespace();

        match words.next() {
            Some("snd") => {
                let word = words.next().unwrap();
                let value = match word.parse() {
                    Ok(value) => value,
                    Err(_)    => *registers.entry(word.clone()).or_insert(0)
                };

                sounds.push(value);

                current += 1;
            },
            Some("set") => {
                let name = words.next().unwrap();
                let word = words.next().unwrap();

                let value : i64 =  match word.parse() {
                    Ok(value) => value,
                    Err(_)    => *registers.entry(word.clone()).or_insert(0)
                };

                registers.insert(name, value);

                current += 1;
            },
            Some("add") => {
                let name = words.next().unwrap();
                let word = words.next().unwrap();

                let value : i64 =  match word.parse() {
                    Ok(value) => value,
                    Err(_)    => *registers.entry(word.clone()).or_insert(0)
                };

                *registers.entry(name.clone()).or_insert(0) += value;

                current += 1;
            },
            Some("mul") => {
                let name = words.next().unwrap();
                let word = words.next().unwrap();

                let value : i64 =  match word.parse() {
                    Ok(value) => value,
                    Err(_)    => *registers.entry(word.clone()).or_insert(0)
                };

                *registers.entry(name.clone()).or_insert(0) *= value;

                current += 1;
            },
            Some("mod") => {
                let name = words.next().unwrap();
                let word = words.next().unwrap();

                let value : i64 =  match word.parse() {
                    Ok(value) => value,
                    Err(_)    => *registers.entry(word.clone()).or_insert(0)
                };

                *registers.entry(name.clone()).or_insert(0) %= value;

                current += 1;
            },
            Some("rcv") => {
                let word = words.next().unwrap();

                let value : i64 =  match word.parse() {
                    Ok(value) => value,
                    Err(_)    => *registers.entry(word.clone()).or_insert(0)
                };

                if value > 0 {
                    return *sounds.iter().last().unwrap();
                }

                current += 1;
            },
            Some("jgz") => {
                let name = words.next().unwrap();
                let word = words.next().unwrap();

                let value : i64 =  match word.parse() {
                    Ok(value) => value,
                    Err(_)    => *registers.entry(word.clone()).or_insert(0)
                };

                if registers[name] > 0 {
                    current = current + value;
                } else {
                    current += 1;
                }
            },

            Some(_) | None => panic!("error")
        }

        println!("-- > current: {}", current);
        println!("-----------------------------------");
    }

    0
}

#[test]
fn execute_test() {
    let input = "set a 1
add a 2
mul a a
mod a 5
snd a
set a 0
rcv a
jgz a -1
set a 1
jgz a -2
";

   assert_eq!(execute(&input), 4);
}
