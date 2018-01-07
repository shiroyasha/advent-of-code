use std::collections::HashMap;

pub struct Program<'a> {
    pub instructions: Vec<&'a str>,
    pub registers: HashMap<&'a str, i64>,
    pub sounds: Vec<i64>,
    pub current: i64,
    pub result: i64
}

impl<'a> Program<'a> {

    pub fn new(source: &'a str) -> Self {
        Program {
            instructions : source.lines().collect(),
            registers: HashMap::new(),
            sounds: Vec::new(),
            current: 0,
            result: 0
        }
    }

    pub fn run(&mut self) {
        while self.current >= 0 && self.current < (self.instructions.len() as i64) {
            for (k, v) in &self.registers {
                print!("{}:{} ", k, v);
            }

            println!("\n{}", self.instructions[self.current as usize]);

            let mut words = self.instructions[self.current as usize].split_whitespace();

            match words.next() {
                Some("snd") => {
                    let word = words.next().unwrap();
                    let value = match word.parse() {
                        Ok(value) => value,
                        Err(_)    => *self.registers.entry(word.clone()).or_insert(0)
                    };

                    self.sounds.push(value);

                    self.current += 1;
                },
                Some("set") => {
                    let name = words.next().unwrap();
                    let word = words.next().unwrap();

                    let value : i64 =  match word.parse() {
                        Ok(value) => value,
                        Err(_)    => *self.registers.entry(word.clone()).or_insert(0)
                    };

                    self.registers.insert(name, value);

                    self.current += 1;
                },
                Some("add") => {
                    let name = words.next().unwrap();
                    let word = words.next().unwrap();

                    let value : i64 =  match word.parse() {
                        Ok(value) => value,
                        Err(_)    => *self.registers.entry(word.clone()).or_insert(0)
                    };

                    *self.registers.entry(name.clone()).or_insert(0) += value;

                    self.current += 1;
                },
                Some("mul") => {
                    let name = words.next().unwrap();
                    let word = words.next().unwrap();

                    let value : i64 =  match word.parse() {
                        Ok(value) => value,
                        Err(_)    => *self.registers.entry(word.clone()).or_insert(0)
                    };

                    *self.registers.entry(name.clone()).or_insert(0) *= value;

                    self.current += 1;
                },
                Some("mod") => {
                    let name = words.next().unwrap();
                    let word = words.next().unwrap();

                    let value : i64 =  match word.parse() {
                        Ok(value) => value,
                        Err(_)    => *self.registers.entry(word.clone()).or_insert(0)
                    };

                    *self.registers.entry(name.clone()).or_insert(0) %= value;

                    self.current += 1;
                },
                Some("rcv") => {
                    let word = words.next().unwrap();

                    let value : i64 =  match word.parse() {
                        Ok(value) => value,
                        Err(_)    => *self.registers.entry(word.clone()).or_insert(0)
                    };

                    if value > 0 {
                        self.result = *self.sounds.iter().last().unwrap();

                        return;
                    }

                    self.current += 1;
                },
                Some("jgz") => {
                    let name = words.next().unwrap();
                    let word = words.next().unwrap();

                    let value : i64 =  match word.parse() {
                        Ok(value) => value,
                        Err(_)    => *self.registers.entry(word.clone()).or_insert(0)
                    };

                    if self.registers[name] > 0 {
                        self.current = self.current + value;
                    } else {
                        self.current += 1;
                    }
                },

                Some(_) | None => panic!("error")
            }

            println!("-- > current: {}", self.current);
            println!("-----------------------------------");
        }
    }

}
