use std::collections::HashMap;

#[derive(PartialEq, Debug)]
pub enum State {
    Pending,
    Running,
    Locked,
    Finished
}

pub struct Program<'a> {
    pub name: String,
    pub instructions: Vec<&'a str>,
    pub registers: HashMap<&'a str, i64>,
    pub sent: Vec<i64>,
    pub incomming: Vec<i64>,
    pub current: i64,
    pub result: i64,
    pub state: State,
    pub sent_count: i64
}

impl<'a> Program<'a> {

    pub fn new(id: i64, source: &'a str) -> Self {
        let mut p = Program {
            name: format!("P{}", id),
            instructions : source.lines().collect(),
            registers: HashMap::new(),
            sent: Vec::new(),
            incomming: Vec::new(),
            state: State::Pending,
            current: 0,
            result: 0,
            sent_count: 0
        };

        p.registers.insert("p", id);

        p
    }

    pub fn run(&mut self) {
        self.state = State::Running;

        while self.current >= 0 && self.current < (self.instructions.len() as i64) {
            let mut words = self.instructions[self.current as usize].split_whitespace();

            print!("{}: {}  ---- ", self.name, self.instructions[self.current as usize]);

            match words.next() {
                Some("snd") => {
                    let word = words.next().unwrap();
                    let value = match word.parse() {
                        Ok(value) => value,
                        Err(_)    => *self.registers.entry(word.clone()).or_insert(0)
                    };

                    self.sent_count += 1;
                    self.sent.push(value);

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

                    if self.incomming.len() == 0 {
                        self.state = State::Locked;

                        return;
                    } else {
                        let value = self.incomming.pop().unwrap();

                        self.registers.insert(word.clone(), value);

                        self.current += 1;
                    }
                },
                Some("jgz") => {
                    let name = words.next().unwrap();
                    let word = words.next().unwrap();

                    let value : i64 = match name.parse() {
                        Ok(v)  => v,
                        Err(_) => self.registers[name]
                    };

                    let offset : i64 =  match word.parse() {
                        Ok(v)  => v,
                        Err(_) => *self.registers.entry(word.clone()).or_insert(0)
                    };

                    if value > 0 {
                        self.current = self.current + offset;
                    } else {
                        self.current += 1;
                    }
                },

                Some(_) | None => panic!("error")
            }

            for (k, v) in &self.registers {
                print!("{}:{} ", k, v);
            }

            println!("");
        }

        self.state = State::Finished;
    }

}
