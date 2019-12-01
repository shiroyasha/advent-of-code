fn spin(steps : usize) -> i64 {
    let mut current_pos : usize = 0;
    let mut res : i64 = 0;

    for i in 1..50_000_000 {
        if i % 100_000 == 0 {
            println!("{}", i);
        }

        current_pos = (current_pos + steps) % i;

        if current_pos == 0 {
            res = i as i64;
        }

        current_pos += 1;
    }

    res
}

#[test]
fn spin_test() {
    assert_eq!(spin(3), 638);
}

fn main() {
    println!("{}", spin(303));
}
