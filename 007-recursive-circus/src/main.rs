// use std::fs::File;
// use std::io::prelude::*;

mod parser;

// #[derive(Debug)]
// struct Program<'a> {
//     name : String,
//     weight: i32,
//     subprogram_names: Vec<String>,
//     programs: Vec<&'a Program<'a>>
// }

// impl<'a> Program<'a> {
//     fn new(name: String, weight: i32, subprogram_names: Vec<String>) -> Program<'a> {
//         Program {
//             name: name,
//             weight: weight,
//             subprogram_names: subprogram_names,
//             programs: vec![]
//         }
//     }

//     fn add(&mut self, program : &Program) {
//         self.programs.push(program);
//     }
// }

// fn parse<'a>(input : &'a str) -> Vec<&Program> {
//     input.lines().map(|line| {
//         let parts : Vec<&str> = line.split_whitespace().collect();

//         let name = parts[0].to_string();
//         let weight = parts[1][1..parts[1].len()-1].parse().unwrap();

//         let subprogram_names : Vec<String> = if parts.len() > 2 {
//             parts[3..].iter().map(|subprogram| {
//                 if subprogram.chars().last().unwrap() == ',' {
//                     subprogram[0..subprogram.len()-1].to_string()
//                 } else {
//                     subprogram[0..].to_string()
//                 }
//             }).collect()
//         } else {
//             vec![]
//         };

//         &Program::new(name, weight, subprogram_names)
//     }).collect()
// }

// fn find_bottom(input: &str) -> String {
//     let mut programs = parse(input);

//     // for p in programs.iter() {
//     //     for subprogram in &p.subprogram_names.iter() {
//     //         p.add(programs.iter().find(|pp| pp.name == subprogram).unwrap().clone());
//     //     }
//     // }

//     // println!("Bottom: {:?}", programs);

//     // sums(&mut map, &programs, &bottom);
//     // let d = diss(&mut map, &programs, &bottom);

//     // println!("Dissbalance: {:?}", map);
//     // println!("Dissbalance: {:?}", d);

//     // bottom
//     "aaa".to_string()
// }

// #[test]
// fn find_bottom_test() {
//     let input = "pbga (66)
// xhth (57)
// ebii (61)
// havc (66)
// ktlj (57)
// fwft (72) -> ktlj, cntj, xhth
// qoyq (66)
// padx (45) -> pbga, havc, qoyq
// tknk (41) -> ugml, padx, fwft
// jptl (61)
// ugml (68) -> gyxo, ebii, jptl
// gyxo (61)
// cntj (57)";

//     assert_eq!(find_bottom(&input), "tknk");
// }

// fn main() {
//     let mut file = File::open("input.txt").expect("Can't open file");
//     let mut content = String::new();

//     file.read_to_string(&mut content).expect("Can't read file");
// }
