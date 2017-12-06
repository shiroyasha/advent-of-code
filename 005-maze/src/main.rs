use std::fs::File;
use std::io::prelude::*;

fn count_jumps(array : &mut Vec<i32>) -> i32 {
    let mut jumps : i32 = 0;
    let mut current_position : usize = 0;
    let len = array.len();

    loop {
        jumps += 1;

        let instruction = array[current_position];
        let next_position : i32 = (current_position as i32) + instruction;

        if next_position >= (len as i32) || next_position < 0 {
            break;
        } else {
            array[current_position] += 1;
        }

        current_position = next_position as usize;
    }

    jumps
}

#[test]
fn count_jumps_test() {
    assert_eq!(count_jumps(&mut vec![0, 3, 0, 1, -3]), 5);
}

fn count_jumps_2(array : &mut Vec<i32>) -> i32 {
    let mut jumps : i32 = 0;
    let mut current_position : usize = 0;
    let len = array.len();

    loop {
        jumps += 1;

        let instruction = array[current_position];
        let next_position : i32 = (current_position as i32) + instruction;

        if next_position >= (len as i32) || next_position < 0 {
            break;
        } else {
            if instruction >= 3 {
                array[current_position] -= 1;
            } else {
                array[current_position] += 1;
            }
        }

        current_position = next_position as usize;
    }

    jumps
}

#[test]
fn count_jumps_test2() {
    assert_eq!(count_jumps_2(&mut vec![0, 3, 0, 1, -3]), 10);
}

fn main() {
    let mut file = File::open("input-part1.txt").expect("Could not open file");
    let mut content = String::new();

    file.read_to_string(&mut content).expect("Could not read file");

    let mut array : Vec<i32> = content.lines().filter_map(|l| l.parse().ok()).collect();
    let mut array2 : Vec<i32> = array.clone();

    let jumps = count_jumps(&mut array);

    println!("Jumps: {}", jumps);

    let jumps2 = count_jumps_2(&mut array2);

    println!("Jumps: {}", jumps2);
}
