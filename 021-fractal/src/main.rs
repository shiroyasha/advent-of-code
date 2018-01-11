use std::fs::File;
use std::io::Read;
use std::fmt;

type Image = Vec<Vec<char>>;

struct Rule {
    input: Image,
    output: Image,
}

fn rotate_image(image: &Image) -> Image {
    let mut result = image.clone();

    for i in 0..image.len() {
        for j in 0..image.len() {
            result[i][j] = image[image.len() - j - 1][i];
        }
    }

    result
}

fn flip_x_image(image: &Image) -> Image {
    let mut result = image.clone();

    for i in 0..image.len() {
        for j in 0..image.len() {
            result[i][j] = image[i][image.len() - j - 1];
        }
    }

    result
}

fn flip_y_image(image: &Image) -> Image {
    let mut result = image.clone();

    for i in 0..image.len() {
        for j in 0..image.len() {
            result[i][j] = image[image.len() - i - 1][j];
        }
    }

    result
}

impl Rule {
    fn matches(&self, image: &Image) -> bool {
        println!("comparing: {:?} to {:?}", self.input, image);

        if image.len() != self.input.len() {
            return false;
        }

        if image.len() != self.input[0].len() {
            return false;
        }

        for i in 0..self.input.len() {
            for j in 0..self.input.len() {
                if image[i][j] != self.input[i][j] {
                    return false;
                }
            }
        }

        true
    }

    fn rotate(&self) -> Rule {
        Rule {
            input: rotate_image(&self.input),
            output: self.output.clone()
        }
    }

    fn flip_x(&self) -> Rule {
        Rule {
            input: flip_x_image(&self.input),
            output: self.output.clone()
        }
    }

    fn flip_y(&self) -> Rule {
        Rule {
            input: flip_y_image(&self.input),
            output: self.output.clone()
        }
    }
}

#[test]
fn matches_test() {
    let image = vec![
        vec!['#', '#'],
        vec!['.', '#'],
    ];

    let rule = Rule {
        input: vec![
            vec!['#', '#'],
            vec!['.', '#'],
        ],
        output: vec![
            vec!['.', '.'],
            vec!['#', '#'],
        ]
    };

    assert!(rule.matches(&image));

    let image = vec![
        vec!['#', '#', '#'],
        vec!['#', '#', '#'],
        vec!['#', '#', '#'],
    ];

    let rule = Rule {
        input: vec![
            vec!['#', '#'],
            vec!['.', '#'],
        ],
        output: vec![
            vec!['.', '.'],
            vec!['#', '#'],
        ]
    };

    assert!(!rule.matches(&image));
}

#[test]
fn rotate_test() {
    let rule = Rule {
        input: vec![
            vec!['#', '#'],
            vec!['.', '#'],
        ],
        output: vec![
            vec!['.', '.', '#'],
            vec!['#', '#', '.'],
            vec!['.', '#', '.'],
        ]
    };

    let rotated = rule.rotate();

    assert_eq!(rotated.input, vec![
       vec!['.', '#'],
       vec!['#', '#'],
    ]);
}

impl fmt::Display for Rule {
    fn fmt(&self, f: &mut fmt::Formatter) -> fmt::Result {
        write!(f, "{:?} => {:?}", self.input, self.output)
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

#[test]
fn count_pixels_test() {
    let image = vec![
        vec!['.', '#', '.'],
        vec!['.', '.', '#'],
        vec!['#', '#', '#'],
    ];

    assert_eq!(count_pixels(&image), 5);
}

fn extract_segment(image: &Image, x: usize, y: usize, w: usize, h: usize) -> Image {
    let mut result = vec![];

    for i in y..(y+h) {
        let mut row = vec![];

        for j in x..(x+w) {
            row.push(image[i][j]);
        }

        result.push(row);
    }

    result
}

#[test]
fn extract_segment_test() {
    let image = vec![
        vec!['.', '#', '.'],
        vec!['.', '.', '#'],
        vec!['#', '#', '#'],
    ];

    assert_eq!(extract_segment(&image, 1, 1, 2, 2), vec![
        vec!['.', '#'],
        vec!['#', '#'],
    ]);
}

fn split_into_segments(image: &Image) -> Vec<Vec<Image>> {
    let mut segments = vec![];

    let segment_size = if image.len() % 2 == 0 { 2 } else { 3 };
    let segment_count = image.len() / segment_size;

    for i in 0..segment_count {
        let mut row = vec![];

        for j in 0..segment_count {
            let image = extract_segment(&image, j * segment_size, i * segment_size, segment_size, segment_size);

            row.push(image);
        }

        segments.push(row);
    }

    segments
}

#[test]
fn split_into_segments_test() {
    let image = vec![
        vec!['.', '#', '.', '#'],
        vec!['.', '.', '#', '.'],
        vec!['#', '#', '#', '#'],
        vec!['#', '.', '.', '#'],
    ];

    assert_eq!(split_into_segments(&image), vec![
        vec![
            vec![
              vec!['.', '#'],
              vec!['.', '.'],
            ],
            vec![
              vec!['.', '#'],
              vec!['#', '.'],
            ],
        ],
        vec![
            vec![
              vec!['#', '#'],
              vec!['#', '.'],
            ],
            vec![
              vec!['#', '#'],
              vec!['.', '#'],
            ],
        ],
    ]);

    let image = vec![
        vec!['.', '#', '.'],
        vec!['.', '.', '#'],
        vec!['#', '#', '#'],
    ];

    assert_eq!(split_into_segments(&image), vec![
        vec![
            vec![
                vec!['.', '#', '.'],
                vec!['.', '.', '#'],
                vec!['#', '#', '#'],
            ],
        ],
    ]);
}

fn join_segments(segments : &Vec<Vec<Image>>) -> Image {
    let mut image = vec![];

    let segment_size = segments[0][0].len();

    let image_width  = segment_size * segments.len();
    let image_height = image_width;

    for i in 0..image_width {
        let mut row = vec![];

        for j in 0..image_height {
            let seg_x = j / segment_size;
            let seg_y = i / segment_size;
            let x = j % segment_size;
            let y = i % segment_size;

            row.push(segments[seg_y][seg_x][y][x]);
        }

        image.push(row);
    }

    image
}

#[test]
fn join_segments_test() {
    let segments = vec![
        vec![
            vec![
                vec!['.', '#'],
                vec!['.', '.'],
            ],
            vec![
                vec!['.', '#'],
                vec!['#', '.'],
            ],
        ],
        vec![
            vec![
                vec!['#', '#'],
                vec!['#', '.'],
            ],
            vec![
                vec!['#', '#'],
                vec!['.', '#'],
            ],
        ],
    ];

    assert_eq!(join_segments(&segments), vec![
        vec!['.', '#', '.', '#'],
        vec!['.', '.', '#', '.'],
        vec!['#', '#', '#', '#'],
        vec!['#', '.', '.', '#'],
    ]);
}

fn process_segment(segment : &Image, rules : &Vec<Rule>) -> Image {
    for rule in rules.iter() {
        let mut flips : Vec<Rule> = vec![];

        flips.push(rule.flip_x().flip_x());
        flips.push(rule.flip_x());
        flips.push(rule.flip_y());
        flips.push(rule.flip_x().flip_y());
        flips.push(rule.flip_y().flip_x());

        for flip in flips.iter() {
            let mut rotations : Vec<Rule> = vec![];

            rotations.push(flip.rotate());
            rotations.push(flip.rotate().rotate());
            rotations.push(flip.rotate().rotate().rotate());
            rotations.push(flip.rotate().rotate().rotate().rotate());

            for rotation in rotations.iter() {
                if rotation.matches(segment) {
                    println!("rule: {}", rotation);
                    return rotation.output.clone();
                }
            }
        }
    }

    panic!("No pattern found");
}

#[test]
fn process_segment_test() {
    let image = vec![
        vec!['#', '#'],
        vec!['.', '#'],
    ];

    let rules = vec![
        Rule {
            input: vec![
                vec!['#', '#'],
                vec!['.', '#'],
            ],
            output: vec![
                vec!['.', '.', '#'],
                vec!['#', '#', '.'],
                vec!['#', '#', '.'],
            ]
        }
    ];

    assert_eq!(process_segment(&image, &rules), vec![
        vec!['.', '.', '#'],
        vec!['#', '#', '.'],
        vec!['#', '#', '.'],
    ]);
}

fn process(image : &Image, rules: &Vec<Rule>) -> Image {
    let segments : Vec<Vec<Image>> = split_into_segments(image).iter().map(|row : &Vec<Image>| {
        row.iter().map(|segment : &Image| process_segment(segment, &rules)).collect()
    }).collect();

    for s in segments.iter() {
        println!("{:?}", s);
    }

    join_segments(&segments)
}

#[test]
fn process_test() {
    let rules = parse("test_input.txt");

    let mut image = vec![
        vec!['.', '#', '.'],
        vec!['.', '.', '#'],
        vec!['#', '#', '#'],
    ];

    for _ in 0..2 {
        image = process(&image, &rules);

        println!("-------------------------------------");
        display(&image);
    }

    assert_eq!(count_pixels(&image), 12);
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

    println!("=====================================");

    display(&image);

    for _ in 0..5 {
        image = process(&image, &rules);

        println!("-------------------------------------");
        display(&image);
    }

    println!("{}", count_pixels(&image));
}
