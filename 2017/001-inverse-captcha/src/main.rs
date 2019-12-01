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
// Input file: part1-input.txt
//
// --- Part Two ---
//
// You notice a progress bar that jumps to 50% completion. Apparently, the door
// isn't yet satisfied, but it did emit a star as encouragement. The instructions
// change:
//
// Now, instead of considering the next digit, it wants you to consider the digit
// halfway around the circular list. That is, if your list contains 10 items, only
// include a digit in your sum if the digit 10/2 = 5 steps forward matches it.
// Fortunately, your list has an even number of elements.
//
// For example:
//
//     1212 produces 6: the list contains 4 items, and all four digits match the
// digit 2 items ahead.  1221 produces 0, because every comparison is between a 1
// and a 2.  123425 produces 4, because both 2s match each other, but no other
// digit has a match.  123123 produces 12.  12131415 produces 4.
//
// What is the solution to your new captcha?
//
// Input file: part2-input.txt
//

use std::fs::File;
use std::io::prelude::*;
use std::iter::*;

// converts characters into vector of numbers
fn to_digits(captcha : String) -> Vec<u32> {
    captcha.chars().filter_map(|s| s.to_digit(10)).collect()
}

// takes [1, 2, 3, 4],
//   returns [(1, 2), (2, 3), (3, 4), (4, 1)] if step is 1
//   returns [(1, 3), (2, 4), (3, 2), (4, 1)] if step is 2
fn each_pair_circular(digits : Vec<u32>, step : usize) -> Vec<(u32, u32)> {
    (0..digits.len())
        .map(|index|               (index, (index+step) % digits.len()))
        .map(|(index, next_index)| (digits[index], digits[next_index]))
        .collect()
}

fn calculate(captcha : String) -> u32 {
    let digits = to_digits(captcha);
    let pairs  = each_pair_circular(digits, 1);

    pairs.iter()
      .filter(|&&(a, b)| a == b) // return pairs that have the same element
      .map(|&(a, _)| a)          // take the first element from array
      .sum()
}

fn calculate2(captcha : String) -> u32 {
    let digits = to_digits(captcha);
    let step   = digits.len() / 2;
    let pairs  = each_pair_circular(digits, step);

    pairs.iter()
      .filter(|&&(a, b)| a == b) // return pairs that have the same element
      .map(|&(a, _)| a)          // take the first element from array
      .sum()
}

fn main() {
    let mut part1_file = File::open("part1_input.txt").expect("File not found");
    let mut part1_input = String::new();
    part1_file.read_to_string(&mut part1_input).expect("Failed to read input file");

    println!("Part 1 Solution: {}", calculate(part1_input));

    let mut part2_file = File::open("part2_input.txt").expect("File not found");
    let mut part2_input = String::new();
    part2_file.read_to_string(&mut part2_input).expect("Failed to read input file");

    println!("Part 2 Solution: {}", calculate2(part2_input));
}

#[cfg(test)]
mod test {
    use super::calculate;
    use super::calculate2;

    #[test]
    fn part1_first_sample_input() {
        assert_eq!(calculate("1122".to_string()), 3);
    }

    #[test]
    fn part1_second_sample_input() {
        assert_eq!(calculate("1111".to_string()), 4);
    }

    #[test]
    fn part1_third_sample_input() {
        assert_eq!(calculate("1234".to_string()), 0);
    }

    #[test]
    fn part1_fourth_sample_input() {
        assert_eq!(calculate("91212129".to_string()), 9);
    }

    #[test]
    fn part1_fifth_sample_input() {
        assert_eq!(calculate("11223311".to_string()), 8);
    }

    #[test]
    fn part2_first_sample() {
        assert_eq!(calculate2("1212".to_string()), 6);
    }

    #[test]
    fn part2_second_sample() {
        assert_eq!(calculate2("1221".to_string()), 0);
    }

    #[test]
    fn part2_third_sample() {
        assert_eq!(calculate2("123425".to_string()), 4);
    }

    #[test]
    fn part2_fourth_sample() {
        assert_eq!(calculate2("123123".to_string()), 12);
    }

    #[test]
    fn part2_fifth_sample() {
        assert_eq!(calculate2("12131415".to_string()), 4);
    }
}
