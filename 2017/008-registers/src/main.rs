use std::collections::HashMap;
use std::fs::File;
use std::io::Read;

struct Registers {
    values: HashMap<String, i32>,
    max_value: i32
}

impl Registers {
    fn new() -> Self {
        Registers { values: HashMap::new(), max_value: 0 }
    }

    fn get(&self, name : &str) -> i32 {
       match self.values.get(name) {
           Some(v) => *v,
           None    => 0
       }
    }

    fn inc(&mut self, name: &str, value: i32) {
        let new_value = self.get(&name) + value;

        if new_value > self.max_value {
            self.max_value = new_value;
        }

        self.values.insert(name.to_string(), new_value);
    }

    fn dec(&mut self, name: &str, value: i32) {
        let new_value = self.get(&name) - value;

        if new_value > self.max_value {
            self.max_value = new_value;
        }

        self.values.insert(name.to_string(), new_value);
    }
}

fn execute(registers : &mut Registers, command : &str) {
    let mut parts = command.split_whitespace();

    let reg_name : &str = parts.next().unwrap();
    let cmd      : &str = parts.next().unwrap();
    let value    : i32  = parts.next().unwrap().parse().unwrap();

    parts.next(); // ignore 'if'

    let test_reg_name : &str = parts.next().unwrap();
    let test_operator : &str = parts.next().unwrap();
    let test_value    : i32  = parts.next().unwrap().parse().unwrap();

    println!("{}, {}, {} - {}, {}, {}", reg_name, cmd, value, test_reg_name, test_operator, test_value);

    let test_result = match test_operator {
        ">"  => registers.get(test_reg_name) > test_value,
        "<"  => registers.get(test_reg_name) < test_value,
        ">=" => registers.get(test_reg_name) >= test_value,
        "<=" => registers.get(test_reg_name) <= test_value,
        "!=" => registers.get(test_reg_name) != test_value,
        "==" => registers.get(test_reg_name) == test_value,
        _    => panic!("Unknown operator")
    };

    if test_result {
        match cmd {
            "inc" => registers.inc(reg_name, value),
            "dec" => registers.dec(reg_name, value),
             _    => panic!("Unknown command")
        };
    }
}

#[test]
fn execute_test() {
    let mut registers = Registers::new();

    execute(&mut registers, "b inc 5 if a > 1");
    execute(&mut registers, "a inc 1 if b < 5");
    execute(&mut registers, "c dec -10 if a >= 1");
    execute(&mut registers, "b inc -20 if c == 10");

    assert_eq!(registers.get("a"), 1);
    assert_eq!(registers.get("b"), -20);
    assert_eq!(registers.get("c"), 10);
}

fn main() {
    let mut content = String::new();
    let mut file = File::open("prog.txt").expect("No such file");

    file.read_to_string(&mut content);

    let mut registers = Registers::new();

    content.lines().for_each(|l| execute(&mut registers, l));

    println!("max values during the operation: {:?}", registers.max_value);
    println!("current max value: {:?}", registers.values.values().max().unwrap());
}
