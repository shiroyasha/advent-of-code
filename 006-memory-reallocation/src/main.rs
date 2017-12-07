
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

fn cycle_count(input : &mut Vec<i32>) -> i32 {
    let mut cycles = 0;
    let mut seen_states : Vec<Vec<i32>> = vec![];

    loop {
        seen_states.push(input.clone());

        cycles += 1;
        reallocate(input);

        println!("{:?}", input);

        if seen_states.iter().filter(|state| state == &input).count() > 0 {
            break;
        }
    }

    cycles
}

#[test]
fn cycle_count_test() {
    let mut input : Vec<i32> = vec![0, 2, 7, 0];

    assert_eq!(cycle_count(&mut input), 5);
}

fn main() {
    let mut input : Vec<i32> = vec![14, 0, 15, 12, 11, 11, 3, 5, 1, 6, 8, 4, 9, 1, 8, 4];

    let cycles = cycle_count(&mut input);
    println!("Cycles {}", cycles);
}
