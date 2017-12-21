fn jugde(input_a : i64, input_b : i64) -> i32 {
    let mut result = 0;

    let mut a = input_a;
    let mut b = input_b;

    for i in 0..40_000_000 {
        a = (a * 16807) % 2147483647;
        b = (b * 48271) % 2147483647;

        // println!("{} : {:0b}", a, (a & 65535));
        // println!("{} : {:0b}", b, (b & 65535));
        // println!("-------------------");

        if (a & 65535) == (b & 65535) {
            result += 1;
        }
    }

    result
}

fn jugde2(input_a : i64, input_b : i64) -> i32 {
    let mut result = 0;

    let mut a = input_a;
    let mut b = input_b;

    let i = 0;

    for _ in 0..5_000_000 {
        loop {
            a = (a * 16807) % 2147483647;

            if a % 4 == 0 { break; }
        }

        loop {
            b = (b * 48271) % 2147483647;

            if b % 8 == 0 { break; }
        }

        if (a & 65535) == (b & 65535) {
            result += 1;
        }
    }

    result
}

#[test]
fn judge_test() {
    assert_eq!(jugde(65, 8921), 588);
    assert_eq!(jugde2(65, 8921), 309);
}
fn main() {
    println!("{}", jugde(618, 814));
    println!("{}", jugde2(618, 814));
}
