
fn reallocate(input : &mut Vec<i32>) {
    let mut index_for_redistribution = 0;
    let mut max_value : i32 = input[0].clone();

    for (index, value) in input.iter().enumerate() {
        if value > &max_value {
            max_value = value.clone();
            index_for_redistribution = index;
        }
    }

    input[index_for_redistribution] = 0;

    while max_value > 0 {
        index_for_redistribution = (index_for_redistribution + 1) % input.len();

        input[index_for_redistribution] += 1;
        max_value -= 1;
    }
}

fn cycle_size(states : &Vec<Vec<i32>>, current_state : &Vec<i32>) -> Option<usize> {
    let mut res = states
        .iter()
        .enumerate()
        .filter(|&(index, state)| state == current_state);

    match res.nth(0) {
        Some((index, _)) => Some(states.len() - index),
        None => None
    }
}

fn cycle_count(input : &mut Vec<i32>) -> (usize, i32) {
    let mut cycles = 0;
    let mut seen_states : Vec<Vec<i32>> = vec![];

    loop {
        seen_states.push(input.clone());

        cycles += 1;
        reallocate(input);

        println!("{:?}", input);

        let cycle_size = cycle_size(&seen_states, &input);

        match cycle_size {
            Some(size) => return (size, cycles),
            None => ()
        }
    }
}

#[test]
fn cycle_count_test() {
    let mut input : Vec<i32> = vec![0, 2, 7, 0];

    let (size, iterations ) = cycle_count(&mut input);

    assert_eq!(iterations, 5);
    assert_eq!(size, 4);
}

fn main() {
    let mut input : Vec<i32> = vec![14, 0, 15, 12, 11, 11, 3, 5, 1, 6, 8, 4, 9, 1, 8, 4];

    let cycles = cycle_count(&mut input);
    println!("Cycles {:?}", cycles);
}
