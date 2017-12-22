use std::fs::File;
use std::io::Read;

fn dance(moves : &str, mut programs : String) -> String {
    moves.split(",").map(|m| m.trim()).for_each(|m| {
        // print!("{} {} => ", programs, m);

        match m.as_bytes()[0] {
            b's' => {
                let value : usize = m[1..].parse().unwrap();
                let split : usize = programs.len() - value;

                programs = format!("{}{}", &programs[split..], &programs[..split]);
            },

            b'x' => {
                let mut parts = m[1..].split("/");

                let pos_a : usize = parts.next().unwrap().parse().unwrap();
                let pos_b : usize = parts.next().unwrap().parse().unwrap();

                let mut program_bytes = programs.clone().into_bytes();

                let tmp = program_bytes[pos_a];

                program_bytes[pos_a] = program_bytes[pos_b];
                program_bytes[pos_b] = tmp;

                programs = String::from_utf8_lossy(&program_bytes).to_string();
            },

            b'p' => {
                let mut parts = m[1..].split("/");

                let pos_a : usize = programs.find(parts.next().unwrap()).unwrap();
                let pos_b : usize = programs.find(parts.next().unwrap()).unwrap();

                let mut program_bytes = programs.clone().into_bytes();

                let tmp = program_bytes[pos_a];

                program_bytes[pos_a] = program_bytes[pos_b];
                program_bytes[pos_b] = tmp;

                programs = String::from_utf8_lossy(&program_bytes).to_string();
            },

            _ => panic!("wtf")
        }

        // println!("{}", programs);
    });

    programs.clone()
}

#[test]
fn dance_test() {
    let moves = "s1,x3/4,pe/b";
    let programs = "abcde";

    assert_eq!(dance(&moves, programs.to_string()), "baedc");
}

fn many_dances(moves : &str, programs : String) {
    let mut solutions : Vec<String> = Vec::new();

    let mut input = programs.clone();
    solutions.push(input.clone());

    //
    // Find the size of the cycle
    //

    let mut i = 0;

    loop {
       let solution = dance(&moves, input.clone());

       println!("{:04}: {} -> {}", i, input, solution);

       match solutions.iter().find(|s| s == &&solution) {
           Some(s) => break,
           None => ()
       }

       solutions.push(solution.clone());

       input = solution;
       i += 1;
    }

    i += 1;

    println!("cycle_size: {}", i);

    // Calculate the 1_000_000_000 iteration by using the power of the %

    let mut input = programs.clone();

    for i in 0..(1_000_000_000 % i) {
       input = dance(&moves, input.clone());
    }

    println!("{}", input);
}

fn main() {
    let mut moves = String::new();
    let mut file = File::open("input.txt").expect("can't open file");
    file.read_to_string(&mut moves);

    let mut programs : String = "abcdefghijklmnop".to_string();

    println!("{}", programs);

    many_dances(&moves, "abcdefghijklmnop".to_string());
}
