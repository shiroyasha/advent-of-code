#[derive(Debug)]
struct Program {
    name : String,
    weight: i32,
    subprograms: Vec<String>,
}

impl Program {

    fn parse_all(input : &str) -> Vec<Self> {
        input.lines().map(|line| Program::parse(line)).collect()
    }

    fn parse(line : &str) -> Self {
        let parts : Vec<&str> = line.split("->").collect();
        let name_and_weight : Vec<&str> = parts[0].split_whitespace().collect();

        let name   : String = name_and_weight[0].to_string();
        let weight : i32    = name_and_weight[1][1..name_and_weight[1].len()-1].parse().unwrap();

        if parts.len() == 2 {
            Program::new(name.clone(), weight, Program::parse_subprogram_names(parts[1]))
        } else {
            Program::new(name.clone(), weight, vec![])
        }
    }

    fn parse_subprogram_names(names : &str) -> Vec<String> {
        names.split(',').map(|subprogram| subprogram.trim().to_string()).collect()
    }

    fn new(name : String, weight : i32, subprograms : Vec<String>) -> Self {
        Program { name, weight, subprograms }
    }
}

impl PartialEq for Program {
    fn eq(&self, other: &Program) -> bool {
        self.name == other.name
    }
}

#[test]
fn parse_test() {
    let input = "pbga (66)\n fwft (72) -> ktlj, cntj, xhth";

    let programs = Program::parse_all(&input);

    assert_eq!(programs[0], Program::new("pbga".to_string(), 66, vec![]));
    assert_eq!(programs[1], Program::new("fwft".to_string(), 72, vec!["ktlj".to_string(), "cntj".to_string(), "xhth".to_string()]));
}
