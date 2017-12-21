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

            row.push(b4);
            row.push(b3);
            row.push(b2);
            row.push(b1);
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

fn display(grid : &Vec<Vec<i32>>) {
    for i in 0..30 {
        for j in 0..30 {
            if grid[i][j] == 0 {
                print!("... ");
            } else {
                print!("{:3} ", grid[i][j]);
            }
        }

        println!("");
    }
}

fn display_visited(grid : &Vec<Vec<bool>>) {
    for i in 0..20 {
        for j in 0..20 {
            if grid[i][j] {
                print!(".");
            } else {
                print!("#");
            }
        }

        println!("");
    }
}

fn dfs(i : usize, j : usize,  mut visited : &mut Vec<Vec<bool>>, mut numbers : &mut Vec<Vec<i32>>, counter : i32) -> bool{
    if visited[i][j] {
        return false;
    }

    visited[i][j] = true;
    numbers[i][j] = counter;

    dfs(i-1, j, &mut visited, &mut numbers, counter);
    dfs(i+1, j, &mut visited, &mut numbers, counter);
    dfs(i, j-1, &mut visited, &mut numbers, counter);
    dfs(i, j+1, &mut visited, &mut numbers, counter);

    true
}

fn regions(input : &str) -> i32 {
    let grid = construct(input);

    display(&grid);

    let mut visited : Vec<Vec<bool>> = vec![];
    let mut numbers : Vec<Vec<i32>> = vec![];
    let mut counter = 1;

    for i in 0..130 {
        let mut v_row = vec![];
        let mut n_row = vec![];

        for j in 0..130 {
            n_row.push(0);

            if i == 0 || j == 0 || i == 129 || j == 129 {
                v_row.push(true);
            } else {
                v_row.push(grid[i-1][j-1] == 0);
            }
        }

        numbers.push(n_row);
        visited.push(v_row);
    }

    display_visited(&visited);

    for i in 1..129 {
        for j in 1..129 {
            if dfs(i, j, &mut visited, &mut numbers, counter) {
                counter += 1;
            }
        }
    }

    display(&numbers);
    // numbers.iter().map(|row| { let r : i32 = *row.iter().max().unwrap(); r }).max().unwrap()
    counter - 1
}

#[test]
fn used_test() {
    assert_eq!(used("flqrgnkx"), 8108);
}

#[test]
fn regions_test() {
    assert_eq!(regions("flqrgnkx"), 1242);
}

fn main() {
    println!("{}", used("nbysizxe"));
    println!("{}", regions("nbysizxe"));
}
