fn spin(steps : usize) -> i32 {
    let mut circular_buffer = Vec::with_capacity(50_000_000);
    circular_buffer.push(0);
    let mut current_pos : usize = 0;

    for i in 0..50_000_000 {
        if i % 100_000 == 0 {
            println!("{}", i);
        }
        // println!("{:?}", circular_buffer);

        current_pos = (current_pos + steps) % circular_buffer.len();

        circular_buffer.insert(current_pos, i + 1);
        current_pos += 1;

        // println!("{:?}, {}", circular_buffer, current_pos);
        // println!("----------------------");
    }

    circular_buffer[1]
}

#[test]
fn spin_test() {
    assert_eq!(spin(3), 638);
}

fn main() {
    println!("{}", spin(303));
}
