fn spin(steps : i32) -> i32 {
    let mut circular_buffer = vec![0];
    let mut current_pos = 0;

    for i in 0..2017 {
        // println!("{:?}", circular_buffer);

        for _ in 0..steps {
            current_pos += 1;

            if current_pos >= circular_buffer.len() {
                current_pos = 0;
            }

            // print!("{},", current_pos);
        }


        circular_buffer.insert(current_pos + 1, i + 1);
        current_pos += 1;

        // println!("{:?}, {}", circular_buffer, current_pos);
        // println!("----------------------");
    }


    // move to location after 2017
    current_pos += 1;

    if current_pos >= circular_buffer.len() {
        current_pos = 0;
    }

    circular_buffer[current_pos]
}

#[test]
fn spin_test() {
    assert_eq!(spin(3), 638);
}

fn main() {
    println!("{}", spin(303));
}
