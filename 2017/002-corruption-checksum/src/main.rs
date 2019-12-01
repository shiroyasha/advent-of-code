use std::iter::*;
use std::fs::File;
use std::io::*;

fn checksum(sheet : String) -> u32 {
    let numbers : Vec<Vec<u32>> = sheet.trim().split("\n").map(|row| {
        row.split(" ").filter_map(|number| number.parse().ok()).collect()
    }).collect();

    numbers.iter().fold(0, |result, row| {
        let max = row.iter().max().unwrap();
        let min = row.iter().min().unwrap();

        result + (max - min)
    })
}

// returns the division result of the only two integers in the
// array that divide eachother
fn dividers(row : &Vec<u32>) -> Option<u32> {
    for i in 0..row.len()-1 {
        for j in (i+1)..row.len() {
            let a = row[i];
            let b = row[j];

            if a % b == 0 {
                return Some(a/b);
            } else if b % a == 0 {
                return Some(b/a);
            }
        }
    }

    None
}

fn checksum2(sheet : String) -> u32 {
    let numbers : Vec<Vec<u32>> = sheet.trim().split("\n").map(|row| {
        row.split(" ").filter_map(|number| number.parse().ok()).collect()
    }).collect();

    numbers.iter().fold(0, |result, row| result + dividers(row).unwrap() )
}

fn main() {
    let mut file = File::open("part1_input.txt").expect("Failed to open file");
    let mut content = String::new();

    file.read_to_string(&mut content).expect("Failed to read file content");

    println!("{}", checksum(content));

    let mut file = File::open("part2_input.txt").expect("Failed to open file");
    let mut content = String::new();

    file.read_to_string(&mut content).expect("Failed to read file content");

    println!("{}", checksum2(content));
}

#[cfg(test)]
mod test {
    use super::checksum;
    use super::checksum2;

    #[test]
    fn part1_example1() {
        let sheet = "5  1 9 5
                     7  5  3
                     2 4 6 8";

        assert_eq!(checksum(sheet.to_string()), 18);
    }

    #[test]
    fn part2_example1() {
        let sheet = " 5 9 2 8
                      9 4 7 3
                      3 8 6 5";

        assert_eq!(checksum2(sheet.to_string()), 9);
    }
}
