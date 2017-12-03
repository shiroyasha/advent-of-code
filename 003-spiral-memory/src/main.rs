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

fn main() {
    println!("Result part 1: {}", access_cost(277678));
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
