fn reverse(array : &mut Vec<i32>, position: i32, size: i32) {
    let mut iterations = size / 2;
    let mut start = position;
    let mut end = (position + size - 1) % (array.len() as i32);

    while iterations > 0 {
        // println!("{}: {:?}, {}, {}", iterations, array, start, end);

        let tmp = array[start as usize];

        array[start as usize] = array[end as usize];
        array[end   as usize] = tmp;

        start = (start + 1 + (array.len() as i32)) % (array.len() as i32);
        end   = (end   - 1 + (array.len() as i32)) % (array.len() as i32);

        iterations -= 1;

        // println!("{}: {:?}, {}, {}", iterations +1, array, start, end);
    }
}

#[test]
fn reverse_test() {
    let mut a = vec![0, 1, 2, 3, 4];

    reverse(&mut a, 0, 3);

    assert_eq!(a, vec![2, 1, 0, 3, 4]);

    let mut b = vec![0, 1, 2, 3, 4];

    reverse(&mut b, 3, 4);

    assert_eq!(b, vec![4, 3, 2, 1, 0]);
}

fn dense(array : &Vec<i32>) -> String {
    let mut result = "".to_string();

    for i in 0..16 {
        let mut d = 0;

        for j in 0..16 {
            d = d ^ array[i*16 + j];
        }

        if d < 16 {
            result = format!("{}0{:x}", result, d);
        } else {
            result = format!("{}{:x}", result, d);
        }
    }

    result
}

fn hash(input: &str) -> String {
    let salt = [17, 31, 73, 47, 23];

    let mut current_pos = 0;
    let mut skip_size = 0;
    let mut array : Vec<i32> = Vec::new();
    let mut lengths = Vec::new();

    for i in 0..256 {
        array.push(i);
    }

    for c in input.bytes() {
        lengths.push(c as i32);
    }

    for s in salt.iter() {
        lengths.push(*s);
    }

    for _ in 0..64 {
        for l in &lengths {
            // println!("{:?}, {}, {}, {}", array, current_pos, skip_size, l);

            reverse(&mut array, current_pos, *l);

            current_pos = (current_pos + l + skip_size) % (array.len() as i32);
            skip_size += 1;
        }
    }

    // println!("{:?}, {}, {}", array, current_pos, skip_size);

    dense(&array)
}

#[test]
fn hash_test() {
    assert_eq!(hash(""), "a2582a3a0e66e6e86e3812dcb672a272");
    assert_eq!(hash("AoC 2017"), "33efeb34ea91902bb2f59c9920caa6cd");
    assert_eq!(hash("1,2,3"), "3efbe78a8d82f29979031a4aa0b16a9d");
    assert_eq!(hash("1,2,4"), "63960835bcdc130f0b66d7ff4f6a5a8e");
}

fn main() {
    let input = "192,69,168,160,78,1,166,28,0,83,198,2,254,255,41,12";

    println!("{:?}", hash(input));
}
