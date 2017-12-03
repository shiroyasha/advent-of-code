// You're standing in a room with "digitization quarantine"
// written in LEDs along one wall. The only door is locked, but
// it includes a small interface. "Restricted Area - Strictly
// No Digitized Users Allowed."
//
// It goes on to explain that you may only leave by solving a
// captcha to prove you're not a human. Apparently, you only
// get one millisecond to solve the captcha: too fast for a
// normal human, but it feels like hours to you.
//
// The captcha requires you to review a sequence of digits
// (your puzzle input) and find the sum of all digits that
// match the next digit in the list. The list is circular, so
// the digit after the last digit is the first digit in the
// list.
//
// For example:
//
//     1122 produces a sum of 3 (1 + 2) because the first digit
// (1) matches the second digit and the third digit (2) matches
// the fourth digit.
//
//     1111 produces 4 because each digit (all 1) matches the next.
//
//     1234 produces 0 because no digit matches the next.
//
//     91212129 produces 9 because the only digit that matches
//              the next one is the last digit, 9.
//
// What is the solution to your captcha?
//
// Input file: input.txt

use std::fs::File;
use std::io::prelude::*;
use std::iter::*;

// converts characters into vector of numbers
fn to_digits(captcha : String) -> Vec<u32> {
    captcha.chars().filter_map(|s| s.to_digit(10)).collect()
}

// takes [1, 2, 3, 4], returns [(1, 2), (2, 3), (3, 4), (4, 1)]
fn each_pair_circular(digits : Vec<u32>) -> Vec<(u32, u32)> {
    (0..digits.len())
        .map(|index|               (index, (index+1) % digits.len()))
        .map(|(index, next_index)| (digits[index], digits[next_index]))
        .collect()
}

fn calculate(captcha : String) -> u32 {
    let digits = to_digits(captcha);
    let pairs  = each_pair_circular(digits);

    pairs.iter()
      .filter(|&&(a, b)| a == b) // return pairs that have the same element
      .map(|&(a, _)| a)          // take the first element from array
      .sum()
}

fn main() {
    let mut file = File::open("input.txt").expect("File not found");

    let mut content = String::new();
    file.read_to_string(&mut content).expect("Failed to read input file");

    println!("Content: {}", content);

    let solution = calculate(content);

    println!("Solution: {}", solution);
}

#[cfg(test)]
mod test {
    use super::calculate;

    #[test]
    fn first_sample_input() {
        assert_eq!(calculate("1122".to_string()), 3);
    }

    #[test]
    fn second_sample_input() {
        assert_eq!(calculate("1111".to_string()), 4);
    }

    #[test]
    fn third_sample_input() {
        assert_eq!(calculate("1234".to_string()), 0);
    }

    #[test]
    fn fourth_sample_input() {
        assert_eq!(calculate("91212129".to_string()), 9);
    }

    #[test]
    fn fifth_sample_input() {
        assert_eq!(calculate("11223311".to_string()), 8);
    }
}
