mod knot;

fn construct(input : &str) -> Vec<Vec<i32>> {
    let mut grid : Vec<Vec<i32>> = Vec::new();

    for i in 0..128 {
        let mut row : Vec<i32> = Vec::new();

        let hash = knot::hash(&format!("{}-{}", input, i));

        for (index, c) in hash.chars().enumerate() {
            let b1 = if i32::from_str_radix(&c.to_string(), 16).unwrap() & (1 << 0) != 0 { 1 } else { 0 };
            let b2 = if i32::from_str_radix(&c.to_string(), 16).unwrap() & (1 << 1) != 0 { 1 } else { 0 };
            let b3 = if i32::from_str_radix(&c.to_string(), 16).unwrap() & (1 << 2) != 0 { 1 } else { 0 };
            let b4 = if i32::from_str_radix(&c.to_string(), 16).unwrap() & (1 << 3) != 0 { 1 } else { 0 };

            row.push(b1);
            row.push(b2);
            row.push(b3);
            row.push(b4);
        }

        grid.push(row);
    }

    grid
}

fn used(input : &str) -> i32 {
    let grid = construct(input);

    grid.iter().map(|row : &Vec<i32>| {
        let sum : i32 = row.iter().sum();

        sum
    }).sum()
}

#[test]
fn used_test() {
    assert_eq!(used("flqrgnkx"), 8108);
}

fn main() {
    println!("{}", used("nbysizxe"));
}
