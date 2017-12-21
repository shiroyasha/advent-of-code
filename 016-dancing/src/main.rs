fn dance(moves : &str, programs : String) -> String {

    moves.split(",").map(|m| m.trim()).for_each(|m| {
        match m.as_bytes()[0] {
            b's' => {
                let split = m[1..].parse().unwrap();

                programs = format!("{}{}", &programs[split..], &programs[..split]);
            },

            b'x' => {
                let mut parts = m[1..].split("/");

                let pos_a : usize = parts.next().unwrap().parse().unwrap();
                let pos_b : usize = parts.next().unwrap().parse().unwrap();

                let mut program_bytes = programs.clone();

                let tmp = program_bytes.as_bytes()[pos_a];

                program_bytes[pos_a] = program_bytes.as_bytes()[pos_b];
                program_bytes[pos_b] = tmp;

                programs = String::from_utf8_lossy(program_bytes).to_string();
            },

            b'p' => {

            },

            _ => panic!("wtf")
        }
    });

    programs.clone()
}

#[test]
fn dance_test() {
    let moves = "s1,x3/4,pe/b";
    let programs = "abcde";

    assert_eq!(dance(&moves, programs.to_string()), "baedc");
}

fn main() {
    println!("Hello, world!");
}
