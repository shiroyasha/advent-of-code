use std::fs::File;
use std::io::Read;
use std::fmt;

type Pattern = Vec<char>;
type Image = Vec<Vec<char>>;

struct Rule {
    input: Vec<Pattern>,
    output: Vec<Pattern>,
}

impl fmt::Display for Rule {
    fn fmt(&self, f: &mut fmt::Formatter) -> fmt::Result {
        write!(f, "{:?} => {:?}", self.input.join(&'/'), self.output.join(&'/'))
    }
}

fn parse(filename : &str) -> Vec<Rule> {
    let mut file = File::open(filename).expect("can't open file");
    let mut content = String::new();

    file.read_to_string(&mut content).expect("can't read from file");

    content.lines().map(|line| {
        let mut parts = line.split("=>");

        let input  = parts.next().unwrap().trim().split("/").map(|p| p.chars().collect()).collect();
        let output = parts.next().unwrap().trim().split("/").map(|p| p.chars().collect()).collect();

        Rule { input, output }
    }).collect()
}

fn display(image : &Image) {
    for row in image {
        println!("{:?}", row);
    }
}

fn count_pixels(image : &Image) -> i64 {
    let mut sum = 0;

    for row in image {
        for pixel in row {
            if pixel == &'#' {
                sum += 1;
            }
        }
    }

    sum
}

fn process_segment(segment : &Image, rules : &Vec<Rule>) -> Image {
    vec![]
}

fn join_segments(segments : &Vec<Vec<Image>>) -> Image {
    vec![]
}

fn split_into_segments(segments : &Image) -> Vec<Vec<Image>> {
    vec![]
}

fn process(image : &Image, rules: &Vec<Rule>) -> Image {
    let segments : Vec<Vec<Image>> = split_into_segments(image).iter().map(|row : &Vec<Image>| {
        row.iter().map(|segment : &Image| process_segment(segment, &rules)).collect()
    }).collect();

    join_segments(&segments)
}

fn main() {
    let rules = parse("input.txt");

    let mut image = vec![
        vec!['.', '#', '.'],
        vec!['.', '.', '#'],
        vec!['#', '#', '#'],
    ];

    for r in rules.iter() {
        println!("{}", r);
    }

    for _ in 1..5 {
        image = process(&image, &rules);

        println!("-------------------------------------");
        display(&image);
    }

    println!("{}", count_pixels(&image));
}
