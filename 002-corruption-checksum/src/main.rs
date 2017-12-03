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

fn main() {
    let mut file = File::open("part1_input.txt").expect("Failed to open file");
    let mut content = String::new();

    file.read_to_string(&mut content).expect("Failed to read file content");

    println!("{}", checksum(content));
}

#[cfg(test)]
mod test {
    use super::checksum;

    #[test]
    fn part1_example1() {
        let sheet = "5  1 9 5
7  5  3
2 4 6 8
";
        assert_eq!(checksum(sheet.to_string()), 18);
    }
}
