fn radius(location : i32) -> i32 {
    let mut result = 0;

    loop {
        if result * result < location {
           result += 1;
        } else {
           break;
        }
    }

    result / 2
}

fn access_cost(location: i32) -> i32 {
    let r = radius(location);
    let up = (2 * r + 1) * (2 * r + 1);
    let down = (2 * r - 1) * (2 * r - 1);
    let len = (up - down) / 4;

    let (x, y) = if down < location && location <= down + len {
        (r, location - down - len/2)
    } else if down + len < location && location <= down + 2 * len {
        (location - down - len - len/2, r)
    } else if down + 2 * len < location && location <= down + 3 * len {
        (-r, location - down - 2*len - len/2)
    } else {
        (location - down - 3*len - len/2, -r)
    };

    x.abs() + y.abs()
}

fn adj_sum(table : &Vec<Vec<i32>>, x : usize, y : usize) -> i32 {
    table[y+1][x-1] + table[y+1][x] + table[y+1][x+1] +
    table[y  ][x-1] +             0 + table[y  ][x+1] +
    table[y-1][x-1] + table[y-1][x] + table[y-1][x+1]
}

fn show_table(table : &Vec<Vec<i32>>) {
    for row in table.iter().rev() {
        println!("{:?}", row);
    }

    println!("------------------------");
}

//
// horible code. kill me please
//
fn access_cost_part2(location: i32) -> i32 {
    let mut table = Vec::with_capacity(20);

    for _ in 0..20 {
        table.push(vec![0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0]);
    }

    table[10][10] = 1;

    let mut x = 10;
    let mut y = 10;
    let mut len = 2;

    for r in (1..5) {
        len = r * 2;
        x += 1;

        for i in 0..len {
            let s = adj_sum(&table, x, y + i);

            if s > location {
                return s;
            }

            table[y+i][x] = s;
        }

        x -= 1;
        y += len - 1;

        for i in 0..len {
            let s = adj_sum(&table, x-i, y);

            if s > location {
                return s;
            }

            table[y][x-i] = s;
        }

        x -= len - 1;
        y -= 1;

        for i in 0..len {
            let s = adj_sum(&table, x, y-i);

            if s > location {
                return s;
            }

            table[y-i][x] = s;
        }

        y -= len - 1;
        x += 1;

        for i in 0..len {
            let s = adj_sum(&table, x+i, y);

            if s > location {
                return s;
            }

            table[y][x+i] = s;
        }

        x += len - 1;
    }

    show_table(&table);

    0
}

fn main() {
    println!("Result part 1: {}", access_cost(277678));
    println!("Result part 2: {}", access_cost_part2(277678));
}

#[cfg(test)]
mod test {
    use super::access_cost;

    #[test]
    fn part1_example1() {
        assert_eq!(access_cost(1), 0)
    }

    #[test]
    fn part1_example2() {
        assert_eq!(access_cost(12), 3)
    }

    #[test]
    fn part1_example3() {
        assert_eq!(access_cost(23), 2)
    }

    #[test]
    fn part1_example4() {
        assert_eq!(access_cost(1024), 31)
    }
}
