use std::io::prelude::*;
use std::fs::File;
use std::collections::HashSet;
use std::iter::FromIterator;

fn has_duplicate_words(string: &str) -> bool {
    let mut words = HashSet::new();

    for word in string.split_whitespace() {
        if words.contains(word) {
            return true;
        } else {
            words.insert(word);
        }
    }

    false
}

fn has_matching_anagrams(string: &str) -> bool {
    let words_with_sorted_chars : Vec<String> = string.split_whitespace().map(|s| {
        let mut sorted_chars: Vec<char> = s.chars().collect();

        sorted_chars.sort();

        String::from_iter(sorted_chars)
    }).collect();

    let mut words = HashSet::new();

    for word in words_with_sorted_chars {
        if words.contains(&word) {
            return true;
        } else {
            words.insert(word.clone());
        }
    }

    false
}

#[test]
fn has_matching_anagrams_test() {
    assert_eq!(has_matching_anagrams("abcde fghij"), false);
    assert_eq!(has_matching_anagrams("abcde xyz ecdab is not valid"), true);
    assert_eq!(has_matching_anagrams("a ab abc abd abf abj"), false);
}

fn main() {
    let mut file = File::open("input-part1.txt").expect("Can't open file");
    let mut content = String::new();

    file.read_to_string(&mut content).expect("Reading from file failed");

    let valid_count = content.lines().filter(|p| { !has_duplicate_words(p) }).count();
    println!("Valid Password count: {}", valid_count);

    let valid_count2 = content.lines().filter(|p| { !has_matching_anagrams(p) }).count();
    println!("Valid Password count: {}", valid_count2);
}
